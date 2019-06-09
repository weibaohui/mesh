package populate

import (
	"fmt"
	"github.com/knative/pkg/apis/istio/v1alpha3"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	"strconv"
	"strings"
)

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