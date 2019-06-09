package routing

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio/controllers/app"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	app.Register(ctx, mContext)
	return nil
}
