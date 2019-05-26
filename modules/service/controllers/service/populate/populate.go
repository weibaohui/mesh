package populate

import (
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate/k8sservice"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate/podcontrollers"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

func Service(service *v1.Service, systemNamespace string, os *objectset.ObjectSet) error {
	k8sservice.Populate(service, systemNamespace, os)
	return podcontrollers.Populate(service, os)

}
