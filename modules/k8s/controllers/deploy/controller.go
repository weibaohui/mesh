package deploy

import (
	"context"
	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mctx *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mctx,
		"k8s-deploy-controller",
		mctx.Apps.Apps().V1().Deployment())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(
			mctx.Apps.Apps().V1().Deployment(),
		).WithRateLimiting(10)

	sh := &serviceHandler{
		namespace:   mctx.Namespace,
		deployCache: mctx.Apps.Apps().V1().Deployment().Cache(),
	}
	// println(sh)
	c.Populator = sh.populate
	return nil
}

type serviceHandler struct {
	namespace   string
	deployCache v12.DeploymentCache
}

func (s *serviceHandler) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {

	// deploy := obj.(*v1.Deployment)
	// if deploy == nil {
	// 	return nil
	// }
	// annotations := deploy.GetAnnotations()
	//
	// inject := annotations[constants.IstioInjector]
	// fmt.Println(inject)
	// if inject == "true" {
	// 	os.Add(deploy)
	// }
	return nil
}
