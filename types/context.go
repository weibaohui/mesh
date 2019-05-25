package types

import (
	"context"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/apiextensions.k8s.io"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/apps"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/extensions"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/networking.istio.io"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/rbac"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type contextKey struct{}

type Context struct {
	Namespace string

	Apps       *apps.Factory
	Core       *core.Factory
	Ext        *apiextensions.Factory
	Extensions *extensions.Factory
	K8s        kubernetes.Interface
	Networking *networking.Factory
	RBAC       *rbac.Factory
	Mesh       *mesh.Factory
	Apply      apply.Apply
}

func From(ctx context.Context) *Context {
	return ctx.Value(contextKey{}).(*Context)
}

func NewContext(namespace string, config *rest.Config) *Context {
	context := &Context{
		Namespace:  namespace,
		Apps:       apps.NewFactoryFromConfigOrDie(config),
		Core:       core.NewFactoryFromConfigOrDie(config),
		Ext:        apiextensions.NewFactoryFromConfigOrDie(config),
		Extensions: extensions.NewFactoryFromConfigOrDie(config),
		Networking: networking.NewFactoryFromConfigOrDie(config),
		RBAC:       rbac.NewFactoryFromConfigOrDie(config),
		K8s:        kubernetes.NewForConfigOrDie(config),
	}

	context.Apply = apply.New(context.K8s.Discovery(), apply.NewClientFactory(config))
	return context
}

func (c *Context) Start(ctx context.Context) error {
	return start.All(ctx, 5,
		c.Apps,
		c.Core,
		c.Extensions,
		c.Ext,
		c.Networking,
		c.RBAC,
	)
}

func BuildContext(ctx context.Context, namespace string, config *rest.Config) (context.Context, *Context) {
	c := NewContext(namespace, config)
	return context.WithValue(ctx, contextKey{}, c), c
}

func Register(f func(context.Context, *Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return f(ctx, From(ctx))
	}
}
