package features

import (
	"context"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/types"

	ntypes "github.com/rancher/mapper"
)

type ControllerRegister func(ctx context.Context, rContext *types.Context) error

type FeatureController struct {
	FeatureName  string
	System       bool
	FeatureSpec  v1.FeatureSpec
	Controllers  []ControllerRegister
	OnStop       func() error
	OnChange     func(*v1.Feature) error
	OnStart      func(*v1.Feature) error
  	registered   bool
}

func (f *FeatureController) Register() error {
	if f.registered {
		return nil
	}
	Register(f)
	return nil
}

func (f *FeatureController) Name() string {
	return f.FeatureName
}

func (f *FeatureController) IsSystem() bool {
	return f.System
}

func (f *FeatureController) Spec() v1.FeatureSpec {
	return f.FeatureSpec
}

func (f *FeatureController) Stop() error {
	if f.OnStop != nil {
		return f.OnStop()
	}

	var errs []error
	return ntypes.NewErrors(errs...)
}

func (f *FeatureController) Changed(feature *v1.Feature) error {
	if f.OnChange != nil {
		if err := f.OnChange(feature); err != nil {
			return err
		}
	}

	return nil
}

func (f *FeatureController) Start(ctx context.Context, feature *v1.Feature) error {


	rContext := types.From(ctx)
	for _, reg := range f.Controllers {
		if err := reg(ctx, rContext); err != nil {
			return err
		}
	}
	// todo: make boot faster
	go func() {
		rContext.Start(ctx)
	}()

	if f.OnStart != nil {
		if err := f.OnStart(feature); err != nil {
			return err
		}
	}

	return nil
}
