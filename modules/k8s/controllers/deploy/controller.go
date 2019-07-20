package deploy

import (
	"context"
	"fmt"

	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/sirupsen/logrus"
	"github.com/weibaohui/mesh/modules/istio/pkg/istio/config"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/utils"
	"github.com/weibaohui/mesh/types"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

func Register(ctx context.Context, mctx *types.Context) error {
	s := DeploymentHandler{
		mctx:        mctx,
		namespace:   mctx.Namespace,
		deployCache: mctx.Apps.Apps().V1().Deployment().Cache(),
		log:         logrus.WithField("name", "k8s-deploy"),
	}
	mctx.Apps.Apps().V1().Deployment().OnChange(ctx, "k8s-deploy-change", s.onChange)
	mctx.Apps.Apps().V1().Deployment().OnRemove(ctx, "k8s-deploy-change", s.onRemove)
	return nil
}

type DeploymentHandler struct {
	mctx        *types.Context
	namespace   string
	deployCache v12.DeploymentCache
	log         *logrus.Entry
}

func (d *DeploymentHandler) onRemove(key string, deploy *v1.Deployment) (*v1.Deployment, error) {
	if deploy == nil {
		return nil, nil
	}
	fmt.Println("deploy onRemove", key, deploy.Name, deploy.Namespace)
	return deploy, nil
}
func (d *DeploymentHandler) onChange(key string, deploy *v1.Deployment) (*v1.Deployment, error) {
	if deploy == nil {
		return nil, nil
	}
	fmt.Println("deploy onChange", key, deploy.Name, deploy.Namespace)
	d.fillTypeMeta(deploy)

	// 检查是否需要注入istio
	annotations := deploy.ObjectMeta.GetAnnotations()
	inject := utils.GetValueFrom(annotations, constants.IstioInjectionEnable)
	if inject == "true" {
		// 纳入管理的
		d.updateAppName(deploy)

		if !d.injected(deploy) {
			injectTemplate, err := d.injectTemplate(deploy)
			injectTemplate.ObjectMeta = deploy.ObjectMeta
			if err != nil {
				d.log.Errorf("injectTemplate err :%s", err.Error())
			}
			deploy = injectTemplate
		}
	}
	err := d.updateDeployment(deploy)
	if err != nil {
		d.log.Errorf("Deployment Update err :%s", err.Error())
	}
	return deploy, nil
}

// 填充TypeMeta
func (d *DeploymentHandler) fillTypeMeta(deploy *v1.Deployment) {
	deploy.TypeMeta.APIVersion = "extensions/v1beta1"
	deploy.TypeMeta.Kind = "Deployment"
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
	object := objects[0]
	deployment := &v1.Deployment{}
	bytes, err := json.Marshal(object)
	err = json.Unmarshal(bytes, deployment)

	return deployment, err
}

func (d *DeploymentHandler) updateAppName(deployment *v1.Deployment) {
	// deploy 自身设置appName
	deployLabels := deployment.GetLabels()
	if !utils.HasLabel(deployLabels, "app") {
		merge := utils.Merge(deployLabels, map[string]string{
			"app": deployment.Name,
		})
		deployment.SetLabels(merge)
	}
	if !utils.HasLabel(deployLabels, "version") {
		merge := utils.Merge(deployLabels, map[string]string{
			"version": "v1",
		})
		deployment.SetLabels(merge)
	}
	// pod 设置appName
	meta := &deployment.Spec.Template.ObjectMeta
	if meta != nil {
		podLabels := meta.GetLabels()
		if !utils.HasLabel(podLabels, "app") {
			merge := utils.Merge(podLabels, map[string]string{
				"app": deployment.Name,
			})
			meta.SetLabels(merge)
		}
		if !utils.HasLabel(podLabels, "version") {
			merge := utils.Merge(podLabels, map[string]string{
				"version": "v1",
			})
			meta.SetLabels(merge)
		}
	}
}
func (d *DeploymentHandler) updateDeployment(deployment *v1.Deployment) (err error) {
	deployment, err = d.mctx.Apps.Apps().V1().Deployment().Update(deployment)
	return err
}
