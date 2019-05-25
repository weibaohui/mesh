package v1

import (
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/genericcondition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	FeatureConditionEnabled = condition.Cond("Enabled")
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Feature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeatureSpec   `json:"spec,omitempty"`
	Status FeatureStatus `json:"status,omitempty"`
}

type FeatureSpec struct {
	Description string   `json:"description,omitempty"`
	Enabled     bool     `json:"enable,omitempty"`
	Requires    []string `json:"features,omitempty"`
}

type FeatureStatus struct {
	EnableOverride *bool                               `json:"enableOverride,omitempty"`
	Conditions     []genericcondition.GenericCondition `json:"conditions,omitempty"`
}
