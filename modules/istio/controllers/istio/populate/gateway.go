package populate

import (
	"fmt"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	"strconv"
	"strings"

	"github.com/knative/pkg/apis/istio/v1alpha3"
	"github.com/rancher/wrangler/pkg/objectset"
)

var (
	supportedProtocol = []v1alpha3.PortProtocol{
		v1alpha3.ProtocolHTTP,
		//v1alpha3.ProtocolTCP,
		//v1alpha3.ProtocolGRPC,
		//v1alpha3.ProtocolHTTP2,
	}
)

func Gateway(systemNamespace string, clusterDomain string, output *objectset.ObjectSet) {
	// Istio Gateway
	gws := v1alpha3.GatewaySpec{
		Selector: map[string]string{
			"istio": constants.IstioGateway,
		},
	}

	// http port
	port, _ := strconv.ParseInt(constants.DefaultHTTPOpenPort, 10, 0)
	for _, protocol := range supportedProtocol {
		gws.Servers = append(gws.Servers, v1alpha3.Server{
			Port: v1alpha3.Port{
				Protocol: protocol,
				Number:   int(port),
				Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(protocol)), port),
			},
			Hosts: []string{"*"},
		})
	}

	httpPort, _ := strconv.ParseInt(constants.DefaultHTTPOpenPort, 10, 0)
	if clusterDomain != "" {
		gws.Servers = append(gws.Servers, v1alpha3.Server{
			Port: v1alpha3.Port{
				Protocol: v1alpha3.ProtocolHTTP,
				Number:   int(httpPort),
				Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(v1alpha3.ProtocolHTTP)), httpPort),
			},
			Hosts: []string{clusterDomain},
		})
	}

	gateway := constructors.NewGateway(systemNamespace, constants.MeshGateway, v1alpha3.Gateway{
		Spec: gws,
	})

	output.Add(gateway)
}
