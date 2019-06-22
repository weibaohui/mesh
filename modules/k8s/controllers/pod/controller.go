package pod

import (
	"context"
	"fmt"
	v13 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/weibaohui/mesh/types"
	corev1 "k8s.io/api/core/v1"
)

func Register(ctx context.Context, mctx *types.Context) error {
	s := PodHandler{
		namespace: mctx.Namespace,
		podCache:  mctx.Core.Core().V1().Pod().Cache(),
	}
	mctx.Core.Core().V1().Pod().OnChange(ctx, "k8s-pod-change", s.onChange)
	return nil
}

type PodHandler struct {
	namespace string
	podCache  v13.PodCache
}

func (d *PodHandler) onChange(key string, pod *corev1.Pod) (*corev1.Pod, error) {
	fmt.Println("pod onChange", key, pod.Name, pod.Namespace)
	return pod, nil
}
