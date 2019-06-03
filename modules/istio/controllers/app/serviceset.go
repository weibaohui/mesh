package app

import (
	"context"
	"fmt"
	"github.com/knative/pkg/apis/istio/v1alpha3"
	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/controllers/service/populate"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sort"
	"strconv"
	"strings"
)

var (
	supportedProtocol = []v1alpha3.PortProtocol{
		v1alpha3.ProtocolHTTP,
		// v1alpha3.ProtocolTCP,
		// v1alpha3.ProtocolGRPC,
		// v1alpha3.ProtocolHTTP2,
	}
)

func Register(ctx context.Context, rContext *types.Context) error {
	fmt.Println("Register app-route-gw ")

	c := stackobject.NewGeneratingController(ctx, rContext, "app-route-gw", rContext.Mesh.Mesh().V1().App())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(
			rContext.Mesh.Mesh().V1().App(),
			rContext.Networking.Networking().V1alpha3().DestinationRule(),
			rContext.Networking.Networking().V1alpha3().VirtualService(),
			rContext.Networking.Networking().V1alpha3().Gateway(),
			rContext.Extensions.Extensions().V1beta1().Ingress()).WithRateLimiting(10)

	sh := &serviceHandler{
		systemNamespace: rContext.Namespace,
		deployCache:     rContext.Apps.Apps().V1().Deployment().Cache(),
		appCache:        rContext.Mesh.Mesh().V1().App().Cache(),
		serviceCache:    rContext.Core.Core().V1().Service().Cache(),
		secretCache:     rContext.Core.Core().V1().Secret().Cache(),
	}

	c.Populator = sh.populate
	fmt.Println("Register app-route-gw ")
	return nil
}

type serviceHandler struct {
	systemNamespace string
	deployCache     v12.DeploymentCache
	appCache        v1.AppCache
	serviceCache    corev1controller.ServiceCache
	secretCache     corev1controller.SecretCache
}

func (s serviceHandler) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {
	fmt.Println("app-route-gw ")

	app := obj.(*meshv1.App)
	if app == nil {
		return nil
	}

	fmt.Println(app.Name)
	fmt.Println(len(app.Spec.Revisions))
	for _, v := range app.Spec.Revisions {
		fmt.Println(v.ServiceName)
		fmt.Println(v.Version)
		fmt.Println(v.Weight)
	}

	if len(app.Spec.Revisions) == 0 {
		return nil
	}

	dr := populate.DestinationRuleForService(app)
	os.Add(dr)

	public := false
	for _, rev := range app.Spec.Revisions {
		if rev.Public {
			public = true
		}
	}
	if !public {
		return nil
	}

	domain := app.Name + "." + app.Namespace + ".oauthd.com"
	gwName := app.Name + "-" + app.Namespace + "-gateway"

	// 域名gateway
	Gateway(app.Namespace, domain, gwName, os)

	// 流量拆分vs

	var dests []populate.Dest
	for _, r := range app.Spec.Revisions {
		dests = append(dests, populate.Dest{
			Host:   app.Name,
			Subset: r.Version,
			Weight: r.Weight,
		})
	}
	sort.Slice(dests, func(i, j int) bool {
		return dests[i].Subset < dests[j].Subset
	})

	var services []*corev1.Service
	for i := len(app.Spec.Revisions) - 1; i >= 0; i-- {
		// requirement, err := labels.NewRequirement("app", "==", []string{app.Name})
		// selector := labels.NewSelector().Add(*requirement)
		// services, err := s.serviceCache.List(app.Namespace, selector)
		service, err := s.serviceCache.Get(app.Namespace, app.Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		} else if errors.IsNotFound(err) {
			continue
		}
		if service == nil {
			return nil
		}
		for _, c := range service.Spec.Ports {
			fmt.Println(c.Name, c.Port, c.Protocol)
		}

		services = append(services, service)
		// deepcopy := deployment.DeepCopy()
		// revVs := populate.VirtualServiceFromSpec(false, s.systemNamespace, app.Name, app.Namespace, nil, deepcopy, dests...)
		// os.Add(revVs)
	}

	vs := VirtualServiceFromService(app.Name, app.Namespace, gwName, domain, services, dests)
	os.Add(vs)
	return nil
}

func VirtualServiceFromService(name, namespace, gateWayName, domain string, services []*corev1.Service, dests []populate.Dest) *v1alpha3.VirtualService {
	vs := constructors.NewVirtualService(namespace, name, v1alpha3.VirtualService{})
	vs.Spec.Gateways = append(vs.Spec.Gateways, gateWayName)

	vs.Spec.Hosts = append(vs.Spec.Hosts, domain)
	vs.Spec.Hosts = append(vs.Spec.Hosts, name)
	vs.Spec.Hosts = append(vs.Spec.Hosts, name+"."+namespace+".cluster.local")
	for _, svc := range services {
		route := httpRoutes(gateWayName, svc, dests)
		for _, r := range route {
			vs.Spec.HTTP = append(vs.Spec.HTTP, r)
		}
	}

	return vs
}

func httpRoutes(gwName string, service *corev1.Service, dests []populate.Dest) []v1alpha3.HTTPRoute {

	var result []v1alpha3.HTTPRoute

	for _, port := range service.Spec.Ports {
		publicPort, route := newRoute(service.Namespace, gwName, port, dests)
		if publicPort != "" {

			result = append(result, route)
		}
	}

	return result
}

type Protocol string

const (
	ProtocolTCP   Protocol = "TCP"
	ProtocolUDP   Protocol = "UDP"
	ProtocolSCTP  Protocol = "SCTP"
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolHTTP2 Protocol = "HTTP2"
	ProtocolGRPC  Protocol = "GRPC"
)

func newRoute(svcNamespace, externalGW string, portBinding corev1.ServicePort, dests []populate.Dest) (string, v1alpha3.HTTPRoute) {
	route := v1alpha3.HTTPRoute{}

	gw := []string{"mesh", externalGW}

	httpPort, _ := strconv.ParseUint(constants.DefaultHTTPOpenPort, 10, 64)
	httpsPort, _ := strconv.ParseUint(constants.DefaultHTTPSOpenPort, 10, 64)
	matches := []v1alpha3.HTTPMatchRequest{
		{
			Port:     uint32(httpPort),
			Gateways: gw,
		},
	}

	matches = append(matches,
		v1alpha3.HTTPMatchRequest{
			Port:     uint32(httpsPort),
			Gateways: gw,
		})

	route.Match = matches

	for _, dest := range dests {

		route.Route = append(route.Route, v1alpha3.HTTPRouteDestination{
			Destination: v1alpha3.Destination{
				Host:   fmt.Sprintf("%s.%s.svc.cluster.local", dest.Host, svcNamespace),
				Subset: dest.Subset,
				Port: v1alpha3.PortSelector{
					Number: uint32(portBinding.Port),
				},
			},
			Weight: dest.Weight,
		})

	}

	sourcePort := httpPort
	if portBinding.Protocol == "https" {
		sourcePort = httpsPort
	}
	return fmt.Sprintf("%v/%s", sourcePort, portBinding.Protocol), route
}
func Gateway(systemNamespace string, clusterDomain, gwName string, output *objectset.ObjectSet) {
	// Istio Gateway
	gws := v1alpha3.GatewaySpec{
		Selector: map[string]string{
			"istio": constants.IstioGateway,
		},
	}

	httpPort, _ := strconv.ParseInt(constants.DefaultHTTPOpenPort, 10, 0)
	if clusterDomain != "" {
		gws.Servers = append(gws.Servers, v1alpha3.Server{
			Port: v1alpha3.Port{
				Protocol: v1alpha3.ProtocolHTTP,
				Number:   int(httpPort),
				Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(v1alpha3.ProtocolHTTP)), httpPort),
			},
			Hosts: []string{clusterDomain, "*." + clusterDomain},
		})
	}

	gateway := constructors.NewGateway(systemNamespace, gwName, v1alpha3.Gateway{
		Spec: gws,
	})

	output.Add(gateway)
}
