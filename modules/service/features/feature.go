package features

import (
	"context"
	"github.com/weibaohui/mesh/modules/service/controllers/appweight"
	"github.com/weibaohui/mesh/modules/service/controllers/externalservice"
	"github.com/weibaohui/mesh/modules/service/controllers/routeset"
	"github.com/weibaohui/mesh/modules/service/controllers/service"
	"github.com/weibaohui/mesh/modules/service/controllers/serviceset"
	"github.com/weibaohui/mesh/modules/service/controllers/servicestatus"
	projectv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/features"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	feature := &features.FeatureController{
		FeatureName: "stack",
		FeatureSpec: projectv1.FeatureSpec{
			Description: "Rio Stack Based UX - required",
			Enabled:     true,
		},
		Controllers: []features.ControllerRegister{
			externalservice.Register,
			routeset.Register,
			service.Register,
			serviceset.Register,
			servicestatus.Register,
			appweight.Register,
		},
	}

	return feature.Register()
}
