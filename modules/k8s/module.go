package k8s

import (
	"context"
	"github.com/weibaohui/mesh/modules/k8s/controllers/deploy"
	"github.com/weibaohui/mesh/modules/k8s/controllers/pod"
	"github.com/weibaohui/mesh/pkg/feature"
	"github.com/weibaohui/mesh/types"
)

func Register(ctx context.Context, mContext *types.Context) error {
	f := &feature.FeatureController{
		FeatureName: "k8s-system",
		Controllers: []feature.ControllerRegister{
			deploy.Register,
			pod.Register,
		},
	}

	return f.Register(ctx, mContext)
}
