package podcontrollers

import (
	"fmt"
	"github.com/rancher/wrangler/pkg/objectset"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constructors"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deployment(service *riov1.Service, cp *controllerParams, os *objectset.ObjectSet) {
	fmt.Println("deployment",service.Name,cp,os)

	dep := constructors.NewDeployment(service.Namespace, service.Name, appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      cp.Labels,
			Annotations: map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cp.Scale.Scale,
			Selector: &metav1.LabelSelector{
				MatchLabels: cp.SelectorLabels,
			},
			Template: cp.PodTemplateSpec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
		},
	})

	if service.SystemSpec != nil && service.SystemSpec.DeploymentStrategy != "" {
		dep.Spec.Strategy.Type = appsv1.DeploymentStrategyType(service.SystemSpec.DeploymentStrategy)
	} else {
		if cp.Scale.Scale > 0 && cp.Scale.BatchSize > 0 {
			dep.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
				MaxSurge:       cp.Scale.MaxSurge,
				MaxUnavailable: cp.Scale.MaxUnavailable,
			}
		}
	}

	os.Add(dep)
 }

func daemonset(service *riov1.Service, cp *controllerParams, os *objectset.ObjectSet) {
	ds := constructors.NewDaemonset(service.Namespace, service.Name, appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      cp.Labels,
			Annotations: map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cp.Scale.Scale,
			Selector: &metav1.LabelSelector{
				MatchLabels: cp.SelectorLabels,
			},
			Template: cp.PodTemplateSpec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
		},
	})

	if service.SystemSpec != nil && service.SystemSpec.DeploymentStrategy != "" {
		ds.Spec.Strategy.Type = appsv1.DeploymentStrategyType(service.SystemSpec.DeploymentStrategy)
	} else {
		if cp.Scale.Scale > 0 && cp.Scale.BatchSize > 0 {
			ds.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
				MaxSurge:       cp.Scale.MaxSurge,
				MaxUnavailable: cp.Scale.MaxUnavailable,
			}
		}
	}

	os.Add(ds)
}
