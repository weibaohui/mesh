package modules

import (
	"context"
	"github.com/weibaohui/mesh/modules/istio"
	"github.com/weibaohui/mesh/modules/service"
	"github.com/weibaohui/mesh/types"

 )

func Register(ctx context.Context, meshContext *types.Context) error {

	if err := istio.Register(ctx, meshContext); err != nil {
		return err
	}
	if err := service.Register(ctx, meshContext); err != nil {
		return err
	}
	return nil
}