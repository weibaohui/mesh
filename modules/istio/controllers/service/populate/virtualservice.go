package populate

import (
	"fmt"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	"github.com/weibaohui/mesh/pkg/services"
	"github.com/weibaohui/mesh/pkg/serviceset"
	"strconv"
	"strings"

	"github.com/knative/pkg/apis/istio/v1alpha3"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/pkg/domains"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	privateGw           = "mesh"
	MeshNameHeader      = "X-Mesh-ServiceName"
	MeshNamespaceHeader = "X-Mesh-Namespace"
	MeshPortHeader      = "X-Mesh-ServicePort"
)

func virtualServices(namespace string, clusterDomain *v1.ClusterDomain, service *v1.Service, os *objectset.ObjectSet) error {
	os.Add(virtualServiceFromService(namespace, clusterDomain, service)...)
	os.Add(gateWay(namespace, clusterDomain, service)...)
	return nil
}

// 给vs 创建匹配的gw
func gateWay(systemNamespace string, clusterDomain *v1.ClusterDomain, service *v1.Service) []runtime.Object {
	var result []runtime.Object

	// Istio Gateway
	gws := v1alpha3.GatewaySpec{
		Selector: map[string]string{
			"istio": constants.IstioGateway,
		},
	}
	httpPort, _ := strconv.ParseInt(constants.DefaultHTTPOpenPort, 10, 0)
	externalHost := domains.GetExternalDomainDot(service.Name, service.Namespace, "oauthd.com")
	gws.Servers = append(gws.Servers, v1alpha3.Server{
		Port: v1alpha3.Port{
			Protocol: v1alpha3.ProtocolHTTP,
			Number:   int(httpPort),
			Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(v1alpha3.ProtocolHTTP)), httpPort),
		},
		Hosts: []string{externalHost},
	})

	gateway := constructors.NewGateway(service.Namespace, service.Name+"-gateway", v1alpha3.Gateway{
		Spec: gws,
	})

	result = append(result, gateway)
	return result
}

func httpRoutes(systemNamespace string, service *v1.Service, dests []Dest) ([]v1alpha3.HTTPRoute, bool) {
	external := false
	var result []v1alpha3.HTTPRoute
	autoscale := false
	if service.Spec.MaxScale != nil && service.Spec.Concurrency != nil && service.Spec.MinScale != nil && *service.Spec.MaxScale != *service.Spec.MinScale {
		autoscale = true
	}
	for _, port := range service.Spec.Ports {
		publicPort, route := newRoute(systemNamespace, domains.GetPublicGateway(systemNamespace), !port.InternalOnly, port, dests, true, autoscale, service)
		if publicPort != "" {
			external = true
			result = append(result, route)
		}
	}

	return result, external
}
func newRoute(systemNamespace, externalGW string, published bool, portBinding v1.ContainerPort, dests []Dest, appendHTTPS bool, autoscale bool, svc *v1.Service) (string, v1alpha3.HTTPRoute) {
	route := v1alpha3.HTTPRoute{}

	if portBinding.Protocol == "" {
		portBinding.Protocol = v1.ProtocolHTTP
	}

	if !isProtocolSupported(portBinding.Protocol) {
		return "", route
	}

	gw := []string{privateGw}
	if published {
		gw = append(gw, externalGW)
	}

	httpPort, _ := strconv.ParseUint(constants.DefaultHTTPOpenPort, 10, 64)
	httpsPort, _ := strconv.ParseUint(constants.DefaultHTTPSOpenPort, 10, 64)
	matches := []v1alpha3.HTTPMatchRequest{
		{
			Port:     uint32(httpPort),
			Gateways: gw,
		},
	}
	if appendHTTPS {
		matches = append(matches,
			v1alpha3.HTTPMatchRequest{
				Port:     uint32(httpsPort),
				Gateways: gw,
			})
	}
	route.Match = matches

	if autoscale {
		if route.Headers == nil {
			route.Headers = &v1alpha3.Headers{
				Request: &v1alpha3.HeaderOperations{
					Add: map[string]string{
						MeshNameHeader:      svc.Name,
						MeshNamespaceHeader: svc.Namespace,
						MeshPortHeader:      strconv.Itoa(int(portBinding.Port)),
					},
				},
			}
		}
		route.Retries = &v1alpha3.HTTPRetry{
			PerTryTimeout: "1m",
			Attempts:      3,
		}
	}

	for _, dest := range dests {
		if autoscale && svc.Status.ObservedScale != nil && *svc.Status.ObservedScale == 0 {
			route.Route = append(route.Route, v1alpha3.HTTPRouteDestination{
				Destination: v1alpha3.Destination{
					Host: fmt.Sprintf("%s.%s.svc.cluster.local", "autoscaler", systemNamespace),
					Port: v1alpha3.PortSelector{
						Number: 80,
					},
				},
				Weight: 100,
			})
		} else {
			route.Route = append(route.Route, v1alpha3.HTTPRouteDestination{
				Destination: v1alpha3.Destination{
					Host:   fmt.Sprintf("%s.%s.svc.cluster.local", dest.Host, svc.Namespace),
					Subset: dest.Subset,
					Port: v1alpha3.PortSelector{
						Number: uint32(portBinding.Port),
					},
				},
				Weight: dest.Weight,
			})
		}
	}

	sourcePort := httpPort
	if portBinding.Protocol == "https" {
		sourcePort = httpsPort
	}
	return fmt.Sprintf("%v/%s", sourcePort, portBinding.Protocol), route
}

type Dest struct {
	Host, Subset string
	Weight       int
}

func DestsForService(namespace, name string, service *serviceset.ServiceSet) []Dest {
	var result []Dest
	for _, rev := range service.Revisions {
		_, ver := services.AppAndVersion(rev)
		weight := rev.Spec.ServiceRevision.Weight
		if rev.Status.WeightOverride != nil {
			weight = *rev.Status.WeightOverride
		}
		result = append(result, Dest{
			Host:   fmt.Sprintf("%s.%s.svc.cluster.local", name, namespace),
			Weight: weight,
			Subset: ver,
		})
	}

	return result
}

func virtualServiceFromService(namespace string, clusterDomain *v1.ClusterDomain, service *v1.Service) []runtime.Object {
	var result []runtime.Object

	// virtual service for each revision
	app, version := services.AppAndVersion(service)
	revVs := VirtualServiceFromSpec(false, namespace, app+"-"+version, service.Namespace, clusterDomain, service, Dest{
		Host:   app,
		Subset: version,
		Weight: 100,
	})
	if revVs != nil {
		result = append(result, revVs)
	}

	return result
}

func VirtualServiceFromSpec(aggregated bool, systemNamespace string, name, namespace string, clusterDomain *v1.ClusterDomain, service *v1.Service, dests ...Dest) *v1alpha3.VirtualService {
	routes, external := httpRoutes(systemNamespace, service, dests)
	if len(routes) == 0 {
		return nil
	}
	//
	//if clusterDomain.Status.ClusterDomain == "" {
	//external = false
	//}

	vs := newVirtualService(name, namespace)

	spec := v1alpha3.VirtualServiceSpec{
		Gateways: []string{privateGw},
		HTTP:     routes,
	}
	if aggregated {
		spec.Hosts = []string{name}
	}

	for _, publicDomain := range service.Status.PublicDomains {
		if publicDomain == "" {
			continue
		}
		spec.Hosts = append(spec.Hosts, publicDomain)
	}

	if external {
		externalGW := domains.GetPublicGateway(systemNamespace)
		spec.Gateways = append(spec.Gateways, externalGW)

		//加入vs自己关联的gateway
		spec.Gateways = append(spec.Gateways, service.Name+"-gateway",)

		externalHost := domains.GetExternalDomainDot(service.Name, service.Namespace, "oauthd.com")
		spec.Hosts = append(spec.Hosts, externalHost)
	}

	vs.Spec = spec
	return vs
}

func newVirtualService(name, namespace string) *v1alpha3.VirtualService {
	return constructors.NewVirtualService(namespace, name, v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
	})
}

func isProtocolSupported(protocol v1.Protocol) bool {
	if protocol == v1.ProtocolHTTP || protocol == v1.ProtocolHTTP2 || protocol == v1.ProtocolGRPC || protocol == v1.ProtocolTCP {
		return true
	}
	return false
}
