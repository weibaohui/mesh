/*

 */

// Code generated by ___go_build_main_go. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=mesh.oauthd.com, Version=v1
	case v1.SchemeGroupVersion.WithResource("apps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Mesh().V1().Apps().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("externalservices"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Mesh().V1().ExternalServices().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("routers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Mesh().V1().Routers().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}