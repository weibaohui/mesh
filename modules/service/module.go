package service

import (
	"context"

	"github.com/weibaohui/mesh/modules/service/features"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	return features.Register(ctx, mContext)
}
