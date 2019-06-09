package feature

import (
	"context"
	"github.com/weibaohui/mesh/types"
)

type FeatureController struct {
	FeatureName string
	Controllers []ControllerRegister
}

func (c *FeatureController) Register(ctx context.Context, mContext *types.Context) error {
	for _, cr := range c.Controllers {
		cr(ctx, mContext)
	}

	return nil
}

type ControllerRegister func(ctx context.Context, mContext *types.Context) error
