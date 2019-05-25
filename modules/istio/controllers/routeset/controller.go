package routeset

import (
	"context"
	"fmt"
	"github.com/weibaohui/mesh/modules/istio/controllers/routeset/populate"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"

	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/rancher/wrangler/pkg/relatedresource"
	"github.com/weibaohui/mesh/modules/istio/pkg/domains"
	adminv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	projectv1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	riov1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	routerDomainUpdate = "router-domain-updater"
)

func Register(ctx context.Context, mContent *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContent, "routing-routeset", mContent.Mesh.Mesh().V1().Router())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(mContent.Networking.Networking().V1alpha3().VirtualService(),
			mContent.Networking.Networking().V1alpha3().DestinationRule(),
			mContent.Networking.Networking().V1alpha3().ServiceEntry(),
			mContent.Extensions.Extensions().V1beta1().Ingress())

	r := &routeSetHandler{
		systemNamespace:      mContent.Namespace,
		secretCache:          mContent.Core.Core().V1().Secret().Cache(),
		externalServiceCache: mContent.Mesh.Mesh().V1().ExternalService().Cache(),
		routesetCache:        mContent.Mesh.Mesh().V1().Router().Cache(),
		clusterDomainCache:   mContent.Mesh.Mesh().V1().ClusterDomain().Cache(),
	}

	mContent.Mesh.Mesh().V1().Router().AddGenericHandler(ctx, routerDomainUpdate, generic.UpdateOnChange(mContent.Mesh.Mesh().V1().Router().Updater(), r.syncDomain))

	relatedresource.Watch(ctx, "externalservice-routeset", r.resolve,
		mContent.Mesh.Mesh().V1().Router(), mContent.Mesh.Mesh().V1().ExternalService())

	c.Populator = r.populate
	return nil
}

func (r routeSetHandler) resolve(namespace, name string, obj runtime.Object) ([]relatedresource.Key, error) {
	switch obj.(type) {
	case *riov1.ExternalService:
		routesets, err := r.routesetCache.List(namespace, labels.Everything())
		if err != nil {
			return nil, err
		}
		var result []relatedresource.Key
		for _, r := range routesets {
			result = append(result, relatedresource.NewKey(r.Namespace, r.Name))
		}
		return result, nil
	}
	return nil, nil
}

type routeSetHandler struct {
	systemNamespace      string
	secretCache          corev1controller.SecretCache
	externalServiceCache riov1controller.ExternalServiceCache
	routesetCache        riov1controller.RouterCache
	clusterDomainCache   projectv1controller.ClusterDomainCache
}

func (r *routeSetHandler) populate(obj runtime.Object, ns *corev1.Namespace, os *objectset.ObjectSet) error {
	routeSet := obj.(*riov1.Router)
	externalServiceMap := map[string]*riov1.ExternalService{}
	routesetMap := map[string]*riov1.Router{}

	clusterDomain, err := r.clusterDomainCache.Get(r.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return err
	}

	ess, err := r.externalServiceCache.List(routeSet.Namespace, labels.Everything())
	if err != nil {
		return err
	}
	for _, es := range ess {
		externalServiceMap[es.Name] = es
	}

	routesets, err := r.routesetCache.List(routeSet.Namespace, labels.Everything())
	if err != nil {
		return err
	}
	for _, rs := range routesets {
		routesetMap[rs.Name] = rs
	}

	if err := populate.VirtualServices(r.systemNamespace, clusterDomain, obj.(*riov1.Router), externalServiceMap, routesetMap, os); err != nil {
		return err
	}

	return nil
}

func (r *routeSetHandler) syncDomain(key string, obj runtime.Object) (runtime.Object, error) {
	if obj == nil {
		return nil, nil
	}

	clusterDomain, err := r.clusterDomainCache.Get(r.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return obj, err
	}

	updateDomain(obj.(*riov1.Router), clusterDomain)

	return obj, nil
}

func updateDomain(router *riov1.Router, clusterDomain *adminv1.ClusterDomain) {
	protocol := "http"
	if clusterDomain.Status.HTTPSSupported {
		protocol = "https"
	}
	router.Status.Endpoints = []string{
		fmt.Sprintf("%s://%s", protocol, domains.GetExternalDomain(router.Name, router.Namespace, clusterDomain.Status.ClusterDomain)),
	}
	for _, pd := range router.Status.PublicDomains {
		router.Status.Endpoints = append(router.Status.Endpoints, fmt.Sprintf("%s://%s", protocol, pd))
	}

	for i, endpoint := range router.Status.Endpoints {
		if protocol == "http" && constants.DefaultHTTPOpenPort != "80" {
			router.Status.Endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPOpenPort)
		}

		if protocol == "https" && constants.DefaultHTTPOpenPort != "443" {
			router.Status.Endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPSOpenPort)
		}
	}
}
