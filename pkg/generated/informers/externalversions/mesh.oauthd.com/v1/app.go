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

// AppInformer provides access to a shared informer and lister for
// Apps.
type AppInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.AppLister
}

type appInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAppInformer constructs a new informer for App type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAppInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAppInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAppInformer constructs a new informer for App type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAppInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MeshV1().Apps(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MeshV1().Apps(namespace).Watch(options)
			},
		},
		&meshoauthdcomv1.App{},
		resyncPeriod,
		indexers,
	)
}

func (f *appInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAppInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *appInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&meshoauthdcomv1.App{}, f.defaultInformer)
}

func (f *appInformer) Lister() v1.AppLister {
	return v1.NewAppLister(f.Informer().GetIndexer())
}
