/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	clientset "github.com/jetstack/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1alpha1"
	informers "github.com/jetstack/cert-manager/pkg/client/informers/externalversions/certmanager/v1alpha1"
	listers "github.com/jetstack/cert-manager/pkg/client/listers/certmanager/v1alpha1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type CertificateHandler func(string, *v1alpha1.Certificate) (*v1alpha1.Certificate, error)

type CertificateController interface {
	CertificateClient

	OnChange(ctx context.Context, name string, sync CertificateHandler)
	OnRemove(ctx context.Context, name string, sync CertificateHandler)
	Enqueue(namespace, name string)

	Cache() CertificateCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type CertificateClient interface {
	Create(*v1alpha1.Certificate) (*v1alpha1.Certificate, error)
	Update(*v1alpha1.Certificate) (*v1alpha1.Certificate, error)
	UpdateStatus(*v1alpha1.Certificate) (*v1alpha1.Certificate, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Certificate, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.CertificateList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Certificate, err error)
}

type CertificateCache interface {
	Get(namespace, name string) (*v1alpha1.Certificate, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.Certificate, error)

	AddIndexer(indexName string, indexer CertificateIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Certificate, error)
}

type CertificateIndexer func(obj *v1alpha1.Certificate) ([]string, error)

type certificateController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.CertificatesGetter
	informer          informers.CertificateInformer
	gvk               schema.GroupVersionKind
}

func NewCertificateController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.CertificatesGetter, informer informers.CertificateInformer) CertificateController {
	return &certificateController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromCertificateHandlerToHandler(sync CertificateHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Certificate
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Certificate))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *certificateController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Certificate))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateCertificateOnChange(updater generic.Updater, handler CertificateHandler) CertificateHandler {
	return func(key string, obj *v1alpha1.Certificate) (*v1alpha1.Certificate, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1alpha1.Certificate)
			}
		}

		return copyObj, err
	}
}

func (c *certificateController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *certificateController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *certificateController) OnChange(ctx context.Context, name string, sync CertificateHandler) {
	c.AddGenericHandler(ctx, name, FromCertificateHandlerToHandler(sync))
}

func (c *certificateController) OnRemove(ctx context.Context, name string, sync CertificateHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromCertificateHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *certificateController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *certificateController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *certificateController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *certificateController) Cache() CertificateCache {
	return &certificateCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *certificateController) Create(obj *v1alpha1.Certificate) (*v1alpha1.Certificate, error) {
	return c.clientGetter.Certificates(obj.Namespace).Create(obj)
}

func (c *certificateController) Update(obj *v1alpha1.Certificate) (*v1alpha1.Certificate, error) {
	return c.clientGetter.Certificates(obj.Namespace).Update(obj)
}

func (c *certificateController) UpdateStatus(obj *v1alpha1.Certificate) (*v1alpha1.Certificate, error) {
	return c.clientGetter.Certificates(obj.Namespace).UpdateStatus(obj)
}

func (c *certificateController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Certificates(namespace).Delete(name, options)
}

func (c *certificateController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Certificate, error) {
	return c.clientGetter.Certificates(namespace).Get(name, options)
}

func (c *certificateController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.CertificateList, error) {
	return c.clientGetter.Certificates(namespace).List(opts)
}

func (c *certificateController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Certificates(namespace).Watch(opts)
}

func (c *certificateController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Certificate, err error) {
	return c.clientGetter.Certificates(namespace).Patch(name, pt, data, subresources...)
}

type certificateCache struct {
	lister  listers.CertificateLister
	indexer cache.Indexer
}

func (c *certificateCache) Get(namespace, name string) (*v1alpha1.Certificate, error) {
	return c.lister.Certificates(namespace).Get(name)
}

func (c *certificateCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.Certificate, error) {
	return c.lister.Certificates(namespace).List(selector)
}

func (c *certificateCache) AddIndexer(indexName string, indexer CertificateIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Certificate))
		},
	}))
}

func (c *certificateCache) GetByIndex(indexName, key string) (result []*v1alpha1.Certificate, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Certificate))
	}
	return result, nil
}
