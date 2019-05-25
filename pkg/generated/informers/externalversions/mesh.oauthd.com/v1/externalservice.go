/*

 */

// Code generated by ___go_build_main_go. DO NOT EDIT.

package v1

import (
	time "time"

	meshoauthdcomv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	versioned "github.com/weibaohui/mesh/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/weibaohui/mesh/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/weibaohui/mesh/pkg/generated/listers/mesh.oauthd.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ExternalServiceInformer provides access to a shared informer and lister for
// ExternalServices.
type ExternalServiceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ExternalServiceLister
}

type externalServiceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewExternalServiceInformer constructs a new informer for ExternalService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewExternalServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredExternalServiceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredExternalServiceInformer constructs a new informer for ExternalService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredExternalServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MeshV1().ExternalServices(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MeshV1().ExternalServices(namespace).Watch(options)
			},
		},
		&meshoauthdcomv1.ExternalService{},
		resyncPeriod,
		indexers,
	)
}

func (f *externalServiceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredExternalServiceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *externalServiceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&meshoauthdcomv1.ExternalService{}, f.defaultInformer)
}

func (f *externalServiceInformer) Lister() v1.ExternalServiceLister {
	return v1.NewExternalServiceLister(f.Informer().GetIndexer())
}