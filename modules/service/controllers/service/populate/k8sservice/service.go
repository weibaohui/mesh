package k8sservice

import (
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/rancher/wrangler/pkg/objectset"
)

func Populate(service *riov1.Service, systemNamespace string, os *objectset.ObjectSet) {
	serviceSelector(service, systemNamespace, os)
	serviceSelector2(service, systemNamespace, os)
 }
