package deploy

import (
	"context"
	"fmt"
	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/weibaohui/mesh/modules/istio/pkg/istio/config"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/utils"
	"github.com/weibaohui/mesh/types"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mctx *types.Context) error {
	s := DeploymentHandler{
		mctx:mctx,
		namespace:   mctx.Namespace,
		deployCache: mctx.Apps.Apps().V1().Deployment().Cache(),
	}
	mctx.Apps.Apps().V1().Deployment().OnChange(ctx, "k8s-deploy-change", s.onChange)
	mctx.Apps.Apps().V1().Deployment().OnRemove(ctx, "k8s-deploy-change", s.onRemove)
	return nil
}

type DeploymentHandler struct {
	mctx        *types.Context
	namespace   string
	deployCache v12.DeploymentCache
}

func (d *DeploymentHandler) onRemove(key string, deploy *v1.Deployment) (*v1.Deployment, error) {
	return deploy,nil
}
func (d *DeploymentHandler) onChange(key string, deploy *v1.Deployment) (*v1.Deployment, error) {
	fmt.Println("deploy onChange", key, deploy.Name, deploy.Namespace)

	annotations := deploy.ObjectMeta.GetAnnotations()
	inject := utils.GetValueFrom(annotations, constants.IstioInjection)
	if inject == "true" {
		if !d.injected(deploy) {
			injectTemplate, err := d.injectTemplate(deploy)
			if err != nil {
				return nil,err
			}
			d.mctx.Apps.Apps().V1().Deployment().Update(injectTemplate)
		}
	}

	return deploy, nil
}

// 是否已经注入过了
func (d *DeploymentHandler) injected(deploy *v1.Deployment) bool {
	for _, c := range deploy.Spec.Template.Spec.Containers {
		if c.Name == constants.IstioProxy {
			return true
		}
	}
	return false
}

func (d *DeploymentHandler) injectTemplate(deploy *v1.Deployment) (*v1.Deployment, error) {
	cm, err := d.mctx.Core.Core().V1().ConfigMap().Get(d.mctx.Namespace, constants.IstionConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	meshConfig, template, err := config.DoConfigAndTemplate(cm.Data[constants.IstioMeshConfigKey], cm.Data[constants.IstioSidecarTemplateName])
	if err != nil {
		return nil, err
	}

	injector := config.NewIstioInjector(meshConfig, template)
	objects, err := injector.Inject([]runtime.Object{deploy})
	if err != nil {
		return nil, err
	}

	utils.YamlToJson(objects)
	return objects[0].(*v1.Deployment), nil
}
