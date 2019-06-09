package stackobject

import (
	"context"
	"fmt"
	"github.com/rancher/mapper/convert"
	"github.com/rancher/wrangler/pkg/condition"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/types"
	"strings"

	"github.com/pkg/errors"
	corev1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/objectset"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

var ErrSkipObjectSet = errors.New("skip objectset")

type ControllerWrapper interface {
	Informer() cache.SharedIndexInformer
	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Enqueue(namespace, name string)
	Updater() generic.Updater
}

type Populator func(obj runtime.Object, ns *corev1.Namespace, os *objectset.ObjectSet) error

type Controller struct {
	Apply          apply.Apply
	Populator      Populator
	name           string
	indexer        cache.Indexer
	namespaceCache corev1controller.NamespaceCache
	injectors      []string
}

func NewGeneratingController(ctx context.Context, mContext *types.Context, name string, controller ControllerWrapper, injectors ...string) *Controller {
	sc := &Controller{
		name:           name,
		Apply:          mContext.Apply.WithSetID(name).WithStrictCaching(),
		namespaceCache: mContext.Core.Core().V1().Namespace().Cache(),
		injectors:      injectors,
		indexer:        controller.Informer().GetIndexer(),
	}

	lcName := name + "-object-controller"
	controller.AddGenericHandler(ctx, lcName, generic.UpdateOnChange(controller.Updater(), sc.OnChange))
	controller.AddGenericRemoveHandler(ctx, lcName+"-remove", sc.OnRemove)
	return sc
}

func (o *Controller) OnRemove(key string, obj runtime.Object) (runtime.Object, error) {
	return obj, o.Apply.WithOwner(obj).Apply(nil)
}

func (o *Controller) OnChange(key string, obj runtime.Object) (runtime.Object, error) {
	if obj == nil {
		return obj, nil
	}

	if o.Populator == nil {
		return obj, nil
	}

	meta, err := meta.Accessor(obj)
	if err != nil {
		return obj, err
	}

	if meta.GetDeletionTimestamp() != nil {
		return obj, err
	}

	ns, err := o.namespaceCache.Get(meta.GetNamespace())
	if err != nil {
		return obj, err
	}

	os := objectset.NewObjectSet()
	if err := o.Populator(obj, ns, os); err != nil {
		if err == ErrSkipObjectSet {
			return obj, nil
		}
		os.AddErr(err)
	}

	desireset := o.Apply.WithOwner(obj)
	if svc, ok := obj.(*v1.Service); ok && !svc.Spec.DisableServiceMesh {
		for _, i := range o.injectors {
			desireset = desireset.WithInjectorName(i)
		}
	}
	return obj, o.getCondition().Do(func() (runtime.Object, error) {
		return obj, desireset.Apply(os)
	})
	return obj, nil

}

func (o *Controller) getCondition() condition.Cond {
	buffer := strings.Builder{}
	buffer.WriteString(string(v1.DeployedCondition))
	for _, part := range strings.Split(o.name, "-") {
		buffer.WriteString(convert.Capitalize(part))
	}
	return condition.Cond(buffer.String())
}
