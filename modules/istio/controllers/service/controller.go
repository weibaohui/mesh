package service

import (
	"context"
	"fmt"

	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/controllers/service/populate"
	"github.com/weibaohui/mesh/modules/istio/pkg/domains"
	adminv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	adminv1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	riov1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	services2 "github.com/weibaohui/mesh/pkg/services"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	serviceDomainUpdate = "service-domain-update"
	appDomainHandler    = "app-domain-update"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext, "routing-service", mContext.Mesh.Mesh().V1().Service())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(
			mContext.Mesh.Mesh().V1().Service(),
			mContext.Networking.Networking().V1alpha3().DestinationRule(),
			mContext.Networking.Networking().V1alpha3().VirtualService(),
			mContext.Extensions.Extensions().V1beta1().Ingress())

	sh := &serviceHandler{
		systemNamespace:      mContext.Namespace,
		serviceClient:        mContext.Mesh.Mesh().V1().Service(),
		serviceCache:         mContext.Mesh.Mesh().V1().Service().Cache(),
		secretCache:          mContext.Core.Core().V1().Secret().Cache(),
		externalServiceCache: mContext.Mesh.Mesh().V1().ExternalService().Cache(),
		clusterDomainCache:   mContext.Mesh.Mesh().V1().ClusterDomain().Cache(),
	}

	mContext.Mesh.Mesh().V1().Service().OnChange(ctx, serviceDomainUpdate, riov1controller.UpdateServiceOnChange(mContext.Mesh.Mesh().V1().Service().Updater(), sh.syncDomain))
	mContext.Mesh.Mesh().V1().App().OnChange(ctx, appDomainHandler, riov1controller.UpdateAppOnChange(mContext.Mesh.Mesh().V1().App().Updater(), sh.syncAppDomain))
	c.Populator = sh.populate
	return nil
}

type serviceHandler struct {
	systemNamespace      string
	serviceClient        v1.ServiceClient
	serviceCache         v1.ServiceCache
	secretCache          corev1controller.SecretCache
	externalServiceCache v1.ExternalServiceCache
	clusterDomainCache   adminv1controller.ClusterDomainCache
}

func (s *serviceHandler) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {


 	service := obj.(*riov1.Service)
	if service.Spec.DisableServiceMesh {
		return nil
	}

	clusterDomain, err := s.clusterDomainCache.Get(s.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return err
	}


	if err := populate.DestinationRulesAndVirtualServices(s.systemNamespace, clusterDomain, service, os); err != nil {
		return err
	}

	return err
}

func (s *serviceHandler) syncDomain(key string, svc *riov1.Service) (*riov1.Service, error) {


	if svc == nil {
		return svc, nil
	}
	if svc.DeletionTimestamp != nil {
		return svc, nil
	}

	clusterDomain, err := s.clusterDomainCache.Get(s.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return svc, err
	}

	updateDomain(svc, clusterDomain)
	return svc, nil
}

func (s *serviceHandler) syncAppDomain(key string, obj *riov1.App) (*riov1.App, error) {
	if obj == nil {
		return obj, nil
	}
	if obj.DeletionTimestamp != nil {
		return obj, nil
	}

	clusterDomain, err := s.clusterDomainCache.Get(s.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return obj, err
	}

	updateAppDomain(obj, clusterDomain)
	return obj, nil
}

func updateAppDomain(app *riov1.App, clusterDomain *adminv1.ClusterDomain) {
	public := true
	for _, svc := range app.Spec.Revisions {
		if !svc.Public {
			public = false
			break
		}
	}
	protocol := "http"
	if clusterDomain.Status.HTTPSSupported {
		protocol = "https"
	}
	var endpoints []string
	if public && clusterDomain.Status.ClusterDomain != "" {
		endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, domains.GetExternalDomain(app.Name, app.Namespace, clusterDomain.Status.ClusterDomain)))
	}
	for _, pd := range app.Status.PublicDomains {
		endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, pd))
	}

	for i, endpoint := range endpoints {
		if protocol == "http" && constants.DefaultHTTPOpenPort != "80" {
			endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPOpenPort)
		}

		if protocol == "https" && constants.DefaultHTTPOpenPort != "443" {
			endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPSOpenPort)
		}
	}
	app.Status.Endpoints = endpoints
}

func updateDomain(service *riov1.Service, clusterDomain *adminv1.ClusterDomain) {
	public := false
	for _, port := range service.Spec.Ports {
		if !port.InternalOnly {
			public = true
			break
		}
	}

	protocol := "http"
	if clusterDomain.Status.HTTPSSupported {
		protocol = "https"
	}

	var endpoints []string
	if public && clusterDomain.Status.ClusterDomain != "" {
		app, version := services2.AppAndVersion(service)
		endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, domains.GetExternalDomain(app+"-"+version, service.Namespace, clusterDomain.Status.ClusterDomain)))
	}

	for _, pd := range service.Status.PublicDomains {
		endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, pd))
	}

	for i, endpoint := range endpoints {
		if protocol == "http" && constants.DefaultHTTPOpenPort != "80" {
			endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPOpenPort)
		}

		if protocol == "https" && constants.DefaultHTTPOpenPort != "443" {
			endpoints[i] = fmt.Sprintf("%s:%s", endpoint, constants.DefaultHTTPSOpenPort)
		}
	}
	service.Status.Endpoints = endpoints
}
