package app

import (
	"context"
	"fmt"
	v12 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/controllers/app/populate"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sort"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext, "app-route-gw", mContext.Mesh.Mesh().V1().App())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(
			mContext.Mesh.Mesh().V1().App(),
			mContext.Networking.Networking().V1alpha3().DestinationRule(),
			mContext.Networking.Networking().V1alpha3().VirtualService(),
			mContext.Networking.Networking().V1alpha3().Gateway(),
			mContext.Extensions.Extensions().V1beta1().Ingress()).WithRateLimiting(10)

	sh := &serviceHandler{
		systemNamespace: mContext.Namespace,
		deployCache:     mContext.Apps.Apps().V1().Deployment().Cache(),
		appCache:        mContext.Mesh.Mesh().V1().App().Cache(),
		serviceCache:    mContext.Core.Core().V1().Service().Cache(),
		secretCache:     mContext.Core.Core().V1().Secret().Cache(),
	}
	c.Populator = sh.populate
	return nil
}

type serviceHandler struct {
	systemNamespace string
	deployCache     v12.DeploymentCache
	appCache        v1.AppCache
	serviceCache    corev1controller.ServiceCache
	secretCache     corev1controller.SecretCache
}

func (s serviceHandler) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {
	// list, e := s.deployCache.List("", labels.Everything())
	// if e != nil {
	// 	fmt.Println(e)
	// }
	// for _, v := range list {
	// 	println(v.Name)
	// }

	app := obj.(*meshv1.App)
	if app == nil {
		return nil
	}

	if len(app.Spec.Revisions) == 0 {
		return nil
	}

	dr := populate.DestinationRuleForService(app)
	os.Add(dr)

	public := false
	for _, rev := range app.Spec.Revisions {
		if rev.Public {
			public = true
		}
	}
	if !public {
		return nil
	}

	domain := app.Name + "." + app.Namespace + ".oauthd.com"
	gwName := app.Name + "-" + app.Namespace + "-gateway"

	// 域名gateway
	populate.Gateway(app.Namespace, domain, gwName, os)

	// 流量拆分vs

	var dests []populate.Dest
	for _, r := range app.Spec.Revisions {
		dests = append(dests, populate.Dest{
			Host:   app.Name,
			Subset: r.Version,
			Weight: r.Weight,
		})
	}
	sort.Slice(dests, func(i, j int) bool {
		return dests[i].Subset < dests[j].Subset
	})

	var services []*corev1.Service
	for i := len(app.Spec.Revisions) - 1; i >= 0; i-- {
		// requirement, err := labels.NewRequirement("app", "==", []string{app.Name})
		// selector := labels.NewSelector().Add(*requirement)
		// services, err := s.serviceCache.List(app.Namespace, selector)
		service, err := s.serviceCache.Get(app.Namespace, app.Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		} else if errors.IsNotFound(err) {
			continue
		}
		if service == nil {
			return nil
		}
		for _, c := range service.Spec.Ports {
			fmt.Println(c.Name, c.Port, c.Protocol)
		}

		services = append(services, service)
		// deepcopy := deployment.DeepCopy()
		// revVs := populate.VirtualServiceFromSpec(false, s.systemNamespace, app.Name, app.Namespace, nil, deepcopy, dests...)
		// os.Add(revVs)
	}

	vs := populate.VirtualServiceFromService(app.Name, app.Namespace, gwName, domain, services, dests)
	os.Add(vs)
	return nil
}
