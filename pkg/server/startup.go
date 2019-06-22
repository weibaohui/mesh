package server

import (
	"context"
	"github.com/weibaohui/mesh/modules"
	"github.com/weibaohui/mesh/pkg/constructors"
	"github.com/weibaohui/mesh/types"

	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/leader"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Crds = append(crd.NonNamespacedTypes(), crd.NamespacedTypes(
	"BuildTemplate.build.knative.dev/v1alpha1",
	"Image.caching.internal.knative.dev/v1alpha1",

	"App.mesh.oauthd.com/v1",
	"ExternalService.mesh.oauthd.com/v1",
	"Router.mesh.oauthd.com/v1",
	"Service.mesh.oauthd.com/v1",
	"Feature.mesh.oauthd.com/v1",
	"ClusterDomain.mesh.oauthd.com/v1",

	"DestinationRule.networking.istio.io/v1alpha3",
	"Gateway.networking.istio.io/v1alpha3",
	"ServiceEntry.networking.istio.io/v1alpha3",
	"VirtualService.networking.istio.io/v1alpha3",

	"adapter.config.istio.io/v1alpha2",
	"attributemanifest.config.istio.io/v1alpha2",
	"EgressRule.config.istio.io/v1alpha2",
	"handler.config.istio.io/v1alpha2",
	"HTTPAPISpecBinding.config.istio.io/v1alpha2",
	"HTTPAPISpec.config.istio.io/v1alpha2",
	"instance.config.istio.io/v1alpha2",
	"kubernetes.config.istio.io/v1alpha2",
	"kubernetesenv.config.istio.io/v1alpha2",
	"logentry.config.istio.io/v1alpha2",
	"metric.config.istio.io/v1alpha2",
	"Policy.authentication.istio.io/v1alpha1",
	"prometheus.config.istio.io/v1alpha2",
	"QuotaSpecBinding.config.istio.io/v1alpha2",
	"QuotaSpec.config.istio.io/v1alpha2",
	"RouteRule.config.istio.io/v1alpha2",
	"rule.config.istio.io/v1alpha2",
	"stdio.config.istio.io/v1alpha2",
	"template.config.istio.io/v1alpha2",
)...)
var globalContext *types.Context
var globalKubeConfig string

func GlobalKubeConfig() string  {
	return globalKubeConfig
}
func GlobalContext() *types.Context {
	if globalContext == nil {
		return nil
	}
  	return globalContext
}
func Startup(ctx context.Context, systemNamespace, kubeConfig string) error {
	globalKubeConfig=kubeConfig
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return err
	}

	if err := Types(ctx, restConfig); err != nil {
		return err
	}

	// 保存下来，共享
	ctx, meshContext := types.BuildContext(ctx, systemNamespace, restConfig)
	globalContext = meshContext


	namespaceClient := meshContext.Core.Core().V1().Namespace()
	if _, err := namespaceClient.Get(systemNamespace, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		ns := constructors.NewNamespace(systemNamespace, v1.Namespace{})
		if _, err := namespaceClient.Create(ns); err != nil {
			return err
		}
	}

	leader.RunOrDie(ctx, systemNamespace, "mesh", meshContext.K8s,
		func(ctx context.Context) {
			runtime.Must(modules.Register(ctx, meshContext))
			runtime.Must(meshContext.Start(ctx))
 			<-ctx.Done()
		})

	return nil
}

func Types(ctx context.Context, config *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(config)
	if err != nil {
		return err
	}

	factory.BatchCreateCRDs(ctx, Crds...)

	return factory.BatchWait()
}
