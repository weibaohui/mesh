package populate

import (
	"fmt"
	"github.com/knative/pkg/apis/istio/v1alpha3"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)


func VirtualServiceFromService(name, namespace, gateWayName, domain string, services []*corev1.Service, dests []Dest) *v1alpha3.VirtualService {
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

func httpRoutes(gwName string, service *corev1.Service, dests []Dest) []v1alpha3.HTTPRoute {

	var result []v1alpha3.HTTPRoute

	for _, port := range service.Spec.Ports {
		publicPort, route := newRoute(service.Namespace, gwName, port, dests)
		if publicPort != "" {

			result = append(result, route)
		}
	}

	return result
}

func newRoute(svcNamespace, externalGW string, portBinding corev1.ServicePort, dests []Dest) (string, v1alpha3.HTTPRoute) {
	route := v1alpha3.HTTPRoute{}

	gw := []string{constants.MeshGateway, externalGW}

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
