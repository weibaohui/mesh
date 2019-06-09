package k8sservice

import (
	"github.com/rancher/wrangler/pkg/objectset"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

func Populate(service *meshv1.Service, systemNamespace string, os *objectset.ObjectSet) {
	serviceSelector(service, systemNamespace, os)
}
