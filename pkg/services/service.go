package services

import (
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
)

func AppAndVersion(service *v1.Service) (string, string) {
	app := service.Spec.App
	version := service.Spec.Version

	if app == "" {
		app = service.Name
	}
	if version == "" {
		version = constants.DefaultServiceVersion
	}

	return app, version
}
