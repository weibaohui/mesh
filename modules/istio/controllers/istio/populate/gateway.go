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

func Gateway(systemNamespace string, clusterDomain string,  output *objectset.ObjectSet) {
	// Istio Gateway
	gws := v1alpha3.GatewaySpec{
		Selector: map[string]string{
			"app": "istio-gateway",
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

	// // https port
	// if clusterDomain != "" {
	// 	httpsPort, _ := strconv.ParseInt(constants.DefaultHTTPSOpenPort, 10, 0)
	// 	gws.Servers = append(gws.Servers, v1alpha3.Server{
	// 		Port: v1alpha3.Port{
	// 			Protocol: v1alpha3.ProtocolHTTPS,
	// 			Number:   int(httpsPort),
	// 			Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(v1alpha3.ProtocolHTTPS)), httpsPort),
	// 		},
	// 		Hosts: []string{clusterDomain},
	// 		TLS: &v1alpha3.TLSOptions{
	// 			Mode:           v1alpha3.TLSModeSimple,
	// 			CredentialName: issuers.RioWildcardCerts,
	// 		},
	// 	})
	// }
	//
	// for _, pd := range publicdomains {
	// 	httpsPort, _ := strconv.ParseInt(constants.DefaultHTTPSOpenPort, 10, 0)
	// 	gws.Servers = append(gws.Servers, v1alpha3.Server{
	// 		Port: v1alpha3.Port{
	// 			Protocol: v1alpha3.ProtocolHTTPS,
	// 			Number:   int(httpsPort),
	// 			Name:     fmt.Sprintf("%v-%v", strings.ToLower(string(v1alpha3.ProtocolHTTPS)), httpsPort),
	// 		},
	// 		Hosts: []string{pd.Spec.DomainName},
	// 		TLS: &v1alpha3.TLSOptions{
	// 			Mode:           v1alpha3.TLSModeSimple,
	// 			CredentialName: pd.Spec.SecretRef.Name,
	// 		},
	// 	})
	// }

	gateway := constructors.NewGateway(systemNamespace, constants.MeshGateway, v1alpha3.Gateway{
		Spec: gws,
	})
	output.Add(gateway)
}
