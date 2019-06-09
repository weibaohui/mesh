package service

import (
	"context"

	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	meshv1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext, "stack-service", mContext.Mesh.Mesh().V1().Service(), "istio-injecter")
	c.Apply = c.Apply.WithCacheTypes(
		mContext.RBAC.Rbac().V1().Role(),
		mContext.RBAC.Rbac().V1().RoleBinding(),
		mContext.RBAC.Rbac().V1().ClusterRole(),
		mContext.RBAC.Rbac().V1().ClusterRoleBinding(),
		mContext.Apps.Apps().V1().Deployment(),
		mContext.Apps.Apps().V1().DaemonSet(),
		mContext.Core.Core().V1().ServiceAccount(),
		mContext.Core.Core().V1().Service(),
		mContext.Core.Core().V1().Secret(),
	).
		WithRateLimiting(5).
		WithStrictCaching()

	sh := &serviceHandler{
		namespace:     mContext.Namespace,
		serviceClient: mContext.Mesh.Mesh().V1().Service(),
		serviceCache:  mContext.Mesh.Mesh().V1().Service().Cache(),
	}

	c.Populator = sh.populate
	return nil
}

type serviceHandler struct {
	namespace     string
	serviceClient meshv1controller.ServiceController
	serviceCache  meshv1controller.ServiceCache
}

func (s *serviceHandler) populate(obj runtime.Object, ns *corev1.Namespace, os *objectset.ObjectSet) error {

	service := obj.(*meshv1.Service)
	if service.Namespace != s.namespace && service.SystemSpec != nil {
		service = service.DeepCopy()
		service.SystemSpec = nil
	}
	return populate.Service(service, s.namespace, os)
}
