package externalservice

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio/controllers/externalservice/populate"
	apiv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/serviceset"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"

	"github.com/rancher/wrangler/pkg/kv"
	"github.com/rancher/wrangler/pkg/objectset"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext,
		"routing-external-service", mContext.Mesh.Mesh().V1().ExternalService())
	c.Apply = c.Apply.WithCacheTypes(mContext.Networking.Networking().V1alpha3().ServiceEntry(),
		mContext.Networking.Networking().V1alpha3().VirtualService())

	p := populator{
		namespace:          mContext.Namespace,
		serviceCache:       mContext.Mesh.Mesh().V1().Service().Cache(),
		clusterDomainCache: mContext.Mesh.Mesh().V1().ClusterDomain().Cache(),
	}

	c.Populator = p.populate
	return nil
}

type populator struct {
	namespace          string
	serviceCache       v1.ServiceCache
	clusterDomainCache v1.ClusterDomainCache
}

func (p populator) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {
	if err := populate.ServiceEntry(obj.(*apiv1.ExternalService), os); err != nil {
		return err
	}

	if obj.(*apiv1.ExternalService).Spec.Service == "" {
		return nil
	}

	targetStackName, targetServiceName := kv.Split(obj.(*apiv1.ExternalService).Spec.Service, "/")
	svc, err := p.serviceCache.Get(targetStackName, targetServiceName)
	if err != nil {
		return err
	}

	serviceSets, err := serviceset.CollectionServices([]*apiv1.Service{svc})
	if err != nil {
		return err
	}

	serviceSet, ok := serviceSets[svc.Name]
	if !ok {
		return err
	}

	clusterDomain, err := p.clusterDomainCache.Get(p.namespace, constants.ClusterDomainName)
	if err != nil {
		return err
	}

	populate.VirtualServiceForExternalService(p.namespace, obj.(*apiv1.ExternalService), serviceSet, clusterDomain, svc, os)
	return nil
}
