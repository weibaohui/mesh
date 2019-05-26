package externalservice

import (
	"context"

	"github.com/weibaohui/mesh/modules/service/controllers/externalservice/populate"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	v12 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	"github.com/rancher/wrangler/pkg/objectset"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext, "stack-external-service", mContext.Mesh.Mesh().V1().ExternalService())
	c.Apply = c.Apply.WithCacheTypes(mContext.Core.Core().V1().Service(),
		mContext.Core.Core().V1().Endpoints(),
		mContext.Networking.Networking().V1alpha3().VirtualService())

	p := populator{
		serviceCache: mContext.Mesh.Mesh().V1().Service().Cache(),
	}

	c.Populator = p.populate
	return nil
}

type populator struct {
	serviceCache v12.ServiceCache
}

func (p populator) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {
	return populate.ServiceForExternalService(obj.(*riov1.ExternalService), namespace, os)
}
