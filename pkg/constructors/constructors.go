package constructors

import (
	"github.com/knative/pkg/apis/istio/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func NewNamespace(name string, obj v1.Namespace) *v1.Namespace {
	obj.APIVersion = "v1"
	obj.Kind = "SystemNamespace"
	obj.Name = name
	return &obj
}

func NewSecret(namespace, name string, obj v1.Secret) *v1.Secret {
	obj.APIVersion = "v1"
	obj.Kind = "Secret"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewServiceAccount(namespace, name string, obj v1.ServiceAccount) *v1.ServiceAccount {
	obj.APIVersion = "v1"
	obj.Kind = "ServiceAccount"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewGateway(namespace, name string, obj v1alpha3.Gateway) *v1alpha3.Gateway {
	obj.APIVersion = "networking.istio.io/v1alpha3"
	obj.Kind = "Gateway"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewDestinationRule(namespace, name string, obj v1alpha3.DestinationRule) *v1alpha3.DestinationRule {
	obj.APIVersion = "networking.istio.io/v1alpha3"
	obj.Kind = "DestinationRule"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewVirtualService(namespace, name string, obj v1alpha3.VirtualService) *v1alpha3.VirtualService {
	obj.APIVersion = "networking.istio.io/v1alpha3"
	obj.Kind = "VirtualService"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewConfigMap(namespace, name string, obj v1.ConfigMap) *v1.ConfigMap {
	obj.APIVersion = "v1"
	obj.Kind = "ConfigMap"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewDeployment(namespace, name string, obj appsv1.Deployment) *appsv1.Deployment {
	obj.APIVersion = "apps/v1"
	obj.Kind = "Deployment"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewDaemonset(namespace, name string, obj appsv1.Deployment) *appsv1.Deployment {
	obj.APIVersion = "apps/v1"
	obj.Kind = "DaemonSet"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewService(namespace, name string, obj v1.Service) *v1.Service {
	obj.APIVersion = "v1"
	obj.Kind = "Service"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewPersistentVolumeClaim(namespace, name string, obj v1.PersistentVolumeClaim) *v1.PersistentVolumeClaim {
	obj.APIVersion = "v1"
	obj.Kind = "PersistentVolumeClaim"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewIngress(namespace, name string, obj v1beta12.Ingress) *v1beta12.Ingress {
	obj.APIVersion = "extensions/v1beta1"
	obj.Kind = "Ingress"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewEndpoints(namespace, name string, obj v1.Endpoints) *v1.Endpoints {
	obj.APIVersion = "v1"
	obj.Kind = "Endpoints"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewCustomResourceDefinition(namespace, name string, obj v1beta1.CustomResourceDefinition) *v1beta1.CustomResourceDefinition {
	obj.APIVersion = "apiextensions.k8s.io/v1beta1"
	obj.Kind = "CustomResourceDefinition"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

func NewServiceEntry(namespace, name string, obj v1alpha3.ServiceEntry) *v1alpha3.ServiceEntry {
	obj.APIVersion = "networking.istio.io/v1alpha3"
	obj.Kind = "ServiceEntry"
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}
