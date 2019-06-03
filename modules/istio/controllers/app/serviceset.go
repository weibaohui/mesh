package app

import (
	"context"
	"fmt"
	"sort"

	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/istio/controllers/service/populate"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	projectv1controller "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	v1 "github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, rContext *types.Context) error {
	fmt.Println("Register app-route-gw ")

	c := stackobject.NewGeneratingController(ctx, rContext, "app-route-gw", rContext.Mesh.Mesh().V1().App())
	c.Apply = c.Apply.WithStrictCaching().
		WithCacheTypes(
			rContext.Mesh.Mesh().V1().App(),
			rContext.Networking.Networking().V1alpha3().DestinationRule(),
			rContext.Networking.Networking().V1alpha3().VirtualService(),
			rContext.Extensions.Extensions().V1beta1().Ingress()).WithRateLimiting(10)

	sh := &serviceHandler{
		systemNamespace:    rContext.Namespace,
		appCache:           rContext.Mesh.Mesh().V1().App().Cache(),
		serviceCache:       rContext.Mesh.Mesh().V1().Service().Cache(),
		secretCache:        rContext.Core.Core().V1().Secret().Cache(),
		clusterDomainCache: rContext.Mesh.Mesh().V1().ClusterDomain().Cache(),
	}

	c.Populator = sh.populate
	fmt.Println("Register app-route-gw ")
	return nil
}

type serviceHandler struct {
	systemNamespace    string
	appCache           v1.AppCache
	serviceCache       v1.ServiceCache
	secretCache        corev1controller.SecretCache
	clusterDomainCache projectv1controller.ClusterDomainCache
}

func (s serviceHandler) populate(obj runtime.Object, namespace *corev1.Namespace, os *objectset.ObjectSet) error {
	fmt.Println("app-route-gw ")

	app := obj.(*meshv1.App)
	if app == nil {
		return nil
	}

	fmt.Println(app.Name)
	fmt.Println(len(app.Spec.Revisions))
	for _, v := range app.Spec.Revisions {
		fmt.Println(v.ServiceName)
		fmt.Println(v.Version)
		fmt.Println(v.Weight)
	}

	clusterDomain, err := s.clusterDomainCache.Get(s.systemNamespace, constants.ClusterDomainName)
	if err != nil {
		return err
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

	var dests []populate.Dest
	for version, rev := range app.Status.RevisionWeight {
		dests = append(dests, populate.Dest{
			Host:   app.Name,
			Subset: version,
			Weight: rev.Weight,
		})
	}
	sort.Slice(dests, func(i, j int) bool {
		return dests[i].Subset < dests[j].Subset
	})

	var revision *meshv1.Service
	for i := len(app.Spec.Revisions) - 1; i >= 0; i-- {
		revision, err = s.serviceCache.Get(app.Namespace, app.Spec.Revisions[i].ServiceName)
		if err != nil && !errors.IsNotFound(err) {
			return err
		} else if errors.IsNotFound(err) {
			continue
		}
		break
	}
	if revision == nil {
		return nil
	}

	deepcopy := revision.DeepCopy()
	deepcopy.Status.PublicDomains = app.Status.PublicDomains
	revVs := populate.VirtualServiceFromSpec(true, s.systemNamespace, app.Name, app.Namespace, clusterDomain, deepcopy, dests...)
	os.Add(revVs)

	return nil
}
