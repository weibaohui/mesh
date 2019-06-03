package routing

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio/controllers/app"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/features"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {

	feature := &features.FeatureController{
		FeatureName: "istio",
		FeatureSpec: v1.FeatureSpec{
			Description: "Service routing using Istio",
			Enabled:     true,
		},

		Controllers: []features.ControllerRegister{
			//externalservice.Register,
			//istio.Register,
			//routeset.Register,
			//service.Register,
			app.Register,
		},
	}

	return feature.Register()
}
