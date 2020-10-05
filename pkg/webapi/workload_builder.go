package webapi

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weibaohui/mesh/pkg/server"
)

type WorkLoadBuilder struct {
	d *v1.Deployment
}

func (w *WorkLoadBuilder) Load(ns, name string) *WorkLoadBuilder {
	mctx := server.GlobalContext()
	deployment, err := mctx.K8s.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	w.d = deployment
	return w
}
func (w *WorkLoadBuilder) Name(name string) *WorkLoadBuilder {
	w.d.ObjectMeta.Name = name
	return w
}

func (w *WorkLoadBuilder) Labels(labels map[string]string) *WorkLoadBuilder {
	w.d.ObjectMeta.Labels = labels
	w.d.Spec.Template.Labels = labels
	w.d.Spec.Selector.MatchLabels = labels
	return w
}
func (w *WorkLoadBuilder) AddLabels(k, v string) *WorkLoadBuilder {
	labels := w.d.ObjectMeta.Labels
	for k, v := range labels {
		fmt.Println(k, v)
	}
	w.d.Spec.Template.Labels = labels
	w.d.Spec.Selector.MatchLabels = labels
	return w
}
