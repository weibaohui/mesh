package service

import (
	"context"
	"github.com/weibaohui/mesh/modules/service/controllers/appweight"
	"github.com/weibaohui/mesh/modules/service/controllers/externalservice"
	"github.com/weibaohui/mesh/modules/service/controllers/routeset"
	"github.com/weibaohui/mesh/modules/service/controllers/service"
	"github.com/weibaohui/mesh/modules/service/controllers/serviceset"
	"github.com/weibaohui/mesh/modules/service/controllers/servicestatus"
	"github.com/weibaohui/mesh/pkg/feature"

	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	f := &feature.FeatureController{
		FeatureName: "k8s-stack",
		Controllers: []feature.ControllerRegister{
			externalservice.Register,
			routeset.Register,
			service.Register,
			serviceset.Register,
			servicestatus.Register,
			appweight.Register,
		},
	}

	return f.Register(ctx, mContext)
}
