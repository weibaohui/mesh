package v1

import (
	"github.com/rancher/wrangler/pkg/genericcondition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec,omitempty"`
	Status AppStatus `json:"status,omitempty"`
}

type AppSpec struct {
	Revisions []Revision `json:"revisions,omitempty"`
}

type Revision struct {
	Public          bool         `json:"public,omitempty"`
	ServiceName     string       `json:"serviceName,omitempty"`
	Version         string       `json:"Version,omitempty"`
	AdjustedWeight  int          `json:"adjustedWeight,omitempty"`
	Weight          int          `json:"weight,omitempty"`
	Scale           int          `json:"scale,omitempty"`
	ScaleStatus     *ScaleStatus `json:"scaleStatus,omitempty"`
	DeploymentReady bool         `json:"deploymentReady,omitempty"`
}

type ServiceObservedWeight struct {
	LastWrite   metav1.Time `json:"lastWrite,omitempty"`
	Weight      int         `json:"weight,omitempty"`
	ServiceName string      `json:"serviceName,omitempty"`
}

type AppStatus struct {
	PublicDomains  []string                            `json:"publicDomains,omitempty"`
	Endpoints      []string                            `json:"endpoints,omitempty"`
	Conditions     []genericcondition.GenericCondition `json:"conditions,omitempty"`
	RevisionWeight map[string]ServiceObservedWeight    `json:"revisionWeight,omitempty"`
}

type ScaleStatus struct {
	// Total number of ready pods targeted by this deployment.
	Ready int `json:"ready,omitempty"`

	// Total number of unavailable pods targeted by this deployment. This is the total number of pods that are still required for the deployment to have 100% available capacity.
	// They may either be pods that are running but not yet available or pods that still have not been created.
	Unavailable int `json:"unavailable,omitempty"`

	// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
	Available int `json:"available,omitempty"`

	// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
	Updated int `json:"updated,omitempty"`
}
