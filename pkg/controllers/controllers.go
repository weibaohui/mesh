package controllers

import (
	"context"
	"github.com/weibaohui/mesh/pkg/controllers/feature"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, rContext *types.Context) error {
	// Controllers
	if err := feature.Register(ctx, rContext); err != nil {
		return err
	}
	return nil
}
