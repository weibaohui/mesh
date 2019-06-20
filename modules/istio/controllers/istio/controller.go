package istio

import (
	"context"
	"github.com/rancher/wrangler/pkg/apply/injectors"
	"github.com/weibaohui/mesh/modules/istio/pkg/istio/config"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func Register(ctx context.Context, mContext *types.Context) error {

	if err := setupConfigMapAndInjectors(ctx, mContext); err != nil {
		return err
	}

	return nil
}


func setupConfigMapAndInjectors(ctx context.Context, mContext *types.Context) error {
	cm, err := mContext.Core.Core().V1().ConfigMap().Get(mContext.Namespace, constants.IstionConfigMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	meshConfig, template, err := config.DoConfigAndTemplate(cm.Data[constants.IstioMeshConfigKey], cm.Data[constants.IstioSidecarTemplateName])
	if err != nil {
		return err
	}

	injector := config.NewIstioInjector(meshConfig, template)
	injectors.Register(constants.IstioInjector, injector.Inject)
	return nil
}
