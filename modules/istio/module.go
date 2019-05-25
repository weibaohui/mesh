package istio

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio/features/routing"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	return routing.Register(ctx, mContext)
}
