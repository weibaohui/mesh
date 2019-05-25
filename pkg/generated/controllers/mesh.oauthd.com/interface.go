/*

 */

// Code generated by ___go_build_main_go. DO NOT EDIT.

package mesh

import (
	"github.com/rancher/wrangler/pkg/generic"
	clientset "github.com/weibaohui/mesh/pkg/generated/clientset/versioned"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	informers "github.com/weibaohui/mesh/pkg/generated/informers/externalversions/mesh.oauthd.com"
)

type Interface interface {
	V1() v1.Interface
}

type group struct {
	controllerManager *generic.ControllerManager
	informers         informers.Interface
	client            clientset.Interface
}

// New returns a new Interface.
func New(controllerManager *generic.ControllerManager, informers informers.Interface,
	client clientset.Interface) Interface {
	return &group{
		controllerManager: controllerManager,
		informers:         informers,
		client:            client,
	}
}

func (g *group) V1() v1.Interface {
	return v1.New(g.controllerManager, g.client.MeshV1(), g.informers.V1())
}