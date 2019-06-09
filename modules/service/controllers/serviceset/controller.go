package serviceset

import (
	"context"
	"sort"

	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	v1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate/serviceports"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constructors"
	meshv1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	services2 "github.com/weibaohui/mesh/pkg/services"
	"github.com/weibaohui/mesh/pkg/serviceset"
	"github.com/weibaohui/mesh/types"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Register(ctx context.Context, mContext *types.Context) error {
	h := handler{
		namespace:      mContext.Namespace,
		apply:          mContext.Apply.WithSetID("serviceset"),
		apps:           mContext.Mesh.Mesh().V1().App(),
		coreservice:    mContext.Core.Core().V1().Service(),
		services:       mContext.Mesh.Mesh().V1().Service(),
		serviceCache:   mContext.Mesh.Mesh().V1().Service().Cache(),
		namespaceCache: mContext.Core.Core().V1().Namespace().Cache(),
	}
	mContext.Mesh.Mesh().V1().Service().OnChange(ctx, "serviceset-controller", h.onChange)
	return nil
}

type handler struct {
	namespace      string
	apply          apply.Apply
	coreservice    corev1controller.ServiceController
	services       meshv1controller.ServiceController
	apps           meshv1controller.AppController
	serviceCache   meshv1controller.ServiceCache
	namespaceCache v1.NamespaceCache
}

func (h *handler) onChange(key string, service *meshv1.Service) (*meshv1.Service, error) {
	os := objectset.NewObjectSet()
	if service == nil {
		return service, nil
	}

	appName, _ := services2.AppAndVersion(service)

	ns, err := h.namespaceCache.Get(service.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return service, nil
		}
		return service, err
	}

	services, err := h.serviceCache.List(service.Namespace, labels.Everything())
	if err != nil {
		return service, err
	}

	serviceSet, err := serviceset.CollectionServices(services)
	if err != nil {
		return service, err
	}
	filteredServices := serviceSet[appName]
	if filteredServices == nil {
		return service, h.apply.WithSetID(appName).
			WithCacheTypes(h.apps).
			WithCacheTypes(h.coreservice).Apply(os)
	}

	svc := createService(ns.Name, appName, filteredServices.Revisions)
	os.Add(svc)

	// ServiceSet
	app := meshv1.NewApp(service.Namespace, appName, meshv1.App{
		Spec: meshv1.AppSpec{
			Revisions: make([]meshv1.Revision, 0),
		},
		Status: meshv1.AppStatus{
			RevisionWeight: make(map[string]meshv1.ServiceObservedWeight, 0),
		},
	})

	var totalweight int
	var serviceWeight []meshv1.Revision
	for _, service := range filteredServices.Revisions {
		if service.DeletionTimestamp != nil {
			continue
		}
		_, version := services2.AppAndVersion(service)
		public := false
		for _, port := range service.Spec.Ports {
			if !port.InternalOnly {
				public = true
				break
			}
		}
		scale := service.Spec.Scale
		if scale == 0 {
			scale = 1
		}
		if service.Status.ObservedScale != nil && *service.Status.ObservedScale != 0 {
			scale = *service.Status.ObservedScale
		}

		scaleStatus := service.Status.ScaleStatus
		weight := service.Spec.Weight

		// hack for daemonsets
		if scaleStatus == nil && service.SystemSpec != nil && service.SystemSpec.Global {
			scaleStatus = &meshv1.ScaleStatus{
				Available: scale,
				Ready:     scale,
			}
			weight = 100
		}

		serviceWeight = append(serviceWeight, meshv1.Revision{
			Public:          public,
			Weight:          weight,
			ServiceName:     service.Name,
			Version:         version,
			Scale:           scale,
			ScaleStatus:     scaleStatus,
			DeploymentReady: IsReady(service.Status.DeploymentStatus),
		})
		totalweight += service.Spec.Weight
	}
	var added int
	for i, rev := range serviceWeight {
		if i == len(serviceWeight)-1 {
			rev.AdjustedWeight = 100 - added
		} else {
			if totalweight == 0 {
				rev.AdjustedWeight = int(1.0 / float64(len(serviceWeight)) * 100)
			} else {
				rev.AdjustedWeight = int(float64(rev.Weight) / float64(totalweight) * 100.0)
			}
			added += rev.AdjustedWeight
		}
		serviceWeight[i] = rev
	}
	sort.Slice(serviceWeight, func(i, j int) bool {
		return serviceWeight[i].Version < serviceWeight[j].Version
	})
	if len(serviceWeight) > 0 {
		app.Spec.Revisions = serviceWeight
		os.Add(app)
	}
	return service, h.apply.WithSetID(appName).
		WithCacheTypes(h.apps).
		WithCacheTypes(h.coreservice).Apply(os)
}

func IsReady(status *appv1.DeploymentStatus) bool {
	if status == nil {
		return false
	}
	for _, con := range status.Conditions {
		if con.Type == appv1.DeploymentAvailable && con.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func createService(namespace, app string, serviceSet []*meshv1.Service) *v12.Service {
	ports := portsForService(serviceSet)
	return constructors.NewService(namespace, app, v12.Service{
		Spec: v12.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": app,
			},
			Type: v12.ServiceTypeClusterIP,
		},
	})
}

func portsForService(serviceSet []*meshv1.Service) (result []v12.ServicePort) {
	ports := map[struct {
		Port     int32
		Protocol v12.Protocol
	}]v12.ServicePort{}

	for _, rev := range serviceSet {
		for _, port := range serviceports.ServiceNamedPorts(rev) {
			ports[struct {
				Port     int32
				Protocol v12.Protocol
			}{
				Port:     port.Port,
				Protocol: port.Protocol,
			}] = port
		}
	}

	for _, port := range ports {
		result = append(result, port)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})

	if len(result) == 0 {
		return []v12.ServicePort{
			{
				Name:       "default",
				Protocol:   v12.ProtocolTCP,
				TargetPort: intstr.FromInt(80),
				Port:       80,
			},
		}
	}
	return
}
