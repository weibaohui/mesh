package serviceset

import (
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/services"
)

func CollectionServices(servicesList []*v1.Service) (Services, error) {
	result := Services{}
	for _, svc := range servicesList {
		app, _ := services.AppAndVersion(svc)

		serviceSet, ok := result[app]
		if !ok {
			serviceSet = &ServiceSet{}
			result[app] = serviceSet
		}
		serviceSet.Revisions = append(serviceSet.Revisions, svc)
	}
	return result, nil
}
