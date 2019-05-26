package routeset

import (
	"context"

	"github.com/weibaohui/mesh/modules/service/controllers/routeset/populate"
	riov1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/stackobject"
	"github.com/weibaohui/mesh/types"
	"github.com/rancher/wrangler/pkg/objectset"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, mContext *types.Context) error {
	c := stackobject.NewGeneratingController(ctx, mContext, "stack-route-set", mContext.Mesh.Mesh().V1().Router())
	c.Apply = c.Apply.WithCacheTypes(mContext.Core.Core().V1().Service(), mContext.Core.Core().V1().Endpoints())

	c.Populator = func(obj runtime.Object, ns *v1.Namespace, os *objectset.ObjectSet) error {
		return populate.ServiceForRouteSet(obj.(*riov1.Router), os)
	}

	return nil
}
