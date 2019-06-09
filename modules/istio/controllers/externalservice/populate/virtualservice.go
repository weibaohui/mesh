package populate

import (
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/controllers/service/populate"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/serviceset"
)

func VirtualServiceForExternalService(namespace string, es *v1.ExternalService, serviceSet *serviceset.ServiceSet, clusterDomain *v1.ClusterDomain,
	svc *v1.Service, os *objectset.ObjectSet) {

	dests := populate.DestsForService(svc.Namespace, svc.Name, serviceSet)
	serviceVS := populate.VirtualServiceFromSpec(true, namespace, svc.Name, svc.Namespace, clusterDomain, svc, dests...)

	// override host match with external service
	serviceVS.Spec.Hosts = []string{}
	serviceVS.Name = es.Name
	serviceVS.Namespace = es.Namespace
	os.Add(serviceVS)
}
