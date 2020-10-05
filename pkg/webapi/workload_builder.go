package webapi

import (
	"k8s.io/kubernetes/pkg/apis/apps"
)

type WorkLoadBuilder struct {
	apps.Deployment
}

func (w *WorkLoadBuilder) Name(name string) *WorkLoadBuilder {
	w.ObjectMeta.Name = name
	return w
}

func (w *WorkLoadBuilder) Labels(labels map[string]string) *WorkLoadBuilder {
	w.ObjectMeta.Labels = labels
	w.Spec.Template.Labels = labels
	w.Spec.Selector.MatchLabels = labels
	return w
}
