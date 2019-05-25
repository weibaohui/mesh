package main

import (
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

var (
	basePackage = "github.com/weibaohui/mesh/types"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/weibaohui/mesh/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"mesh.oauthd.com": {
				Types: []interface{}{
					meshv1.ExternalService{},
					meshv1.Router{},
					meshv1.App{},
					meshv1.Feature{},
					meshv1.Service{},
					meshv1.ClusterDomain{},
				},
				GenerateTypes: true,
			},
		},
	})
}
