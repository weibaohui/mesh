package populate

import (
	"fmt"
	"github.com/weibaohui/mesh/pkg/constructors"

	"github.com/knative/pkg/apis/istio/v1alpha3"
 	"github.com/rancher/wrangler/pkg/objectset"
	apiv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

func DestinationRulesAndVirtualServices(namespace string, clusterDomain *apiv1.ClusterDomain, service *apiv1.Service, os *objectset.ObjectSet) error {
	return virtualServices(namespace, clusterDomain, service, os)
}

func DestinationRuleForService(app *apiv1.App) *v1alpha3.DestinationRule {
	drSpec := v1alpha3.DestinationRuleSpec{
		Host: fmt.Sprintf("%s.%s.svc.cluster.local", app.Name, app.Namespace),
	}

	for _, rev := range app.Spec.Revisions {
		drSpec.Subsets = append(drSpec.Subsets, newSubSet(rev.Version))
	}

	dr := newDestinationRule(app.Namespace, app.Name)
	dr.Spec = drSpec

	return dr
}

func newSubSet(version string) v1alpha3.Subset {
	return v1alpha3.Subset{
		Name: version,
		Labels: map[string]string{
			"version": version,
		},
	}
}

func newDestinationRule(namespace, name string) *v1alpha3.DestinationRule {
	return constructors.NewDestinationRule(namespace, name, v1alpha3.DestinationRule{})
}
