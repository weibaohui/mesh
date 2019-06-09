package k8sservice

import (
	"github.com/rancher/wrangler/pkg/objectset"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

func Populate(service *riov1.Service, systemNamespace string, os *objectset.ObjectSet) {
	serviceSelector(service, systemNamespace, os)
}
