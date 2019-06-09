package istio

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio/controllers/app"
	"github.com/weibaohui/mesh/pkg/feature"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	f := &feature.FeatureController{
		FeatureName: "Istio-stack",
		Controllers: []feature.ControllerRegister{
			app.Register,
		},
	}

	return f.Register(ctx, mContext)
}
