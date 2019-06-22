package deploy

import (
	"context"
	"fmt"
	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/weibaohui/mesh/types"
 	v1 "k8s.io/api/apps/v1"
)

func Register(ctx context.Context, mctx *types.Context) error {
	s := DeploymentHandler{
		namespace:   mctx.Namespace,
		deployCache: mctx.Apps.Apps().V1().Deployment().Cache(),
	}
	mctx.Apps.Apps().V1().Deployment().OnChange(ctx, "k8s-deploy-change",s.onChange)
	return nil
}

type DeploymentHandler struct {
	namespace   string
	deployCache v12.DeploymentCache
}

func (d *DeploymentHandler ) onChange(key string, deploy *v1.Deployment) (*v1.Deployment, error){
	fmt.Println("deploy onChange",key,deploy.Name,deploy.Namespace)
	return deploy, nil
}
