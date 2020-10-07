package webapi

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	labels[k] = v
	w.d.Spec.Template.Labels = labels
	w.d.Spec.Selector.MatchLabels = labels
	return w
}

func (w *WorkLoadBuilder) String() *WorkLoadBuilder {
	bytes, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(bytes))
	return w
}
func (w *WorkLoadBuilder) Containers() []corev1.Container {
	return w.d.Spec.Template.Spec.Containers
}
func (w *WorkLoadBuilder) SetImageByContainerName(image, containerName string) *WorkLoadBuilder {
	containers := w.d.Spec.Template.Spec.Containers
	for i := range containers {
		c := containers[i]
		if c.Name == containerName {
			c.Image = image
			break
		}
	}
	return w
}
