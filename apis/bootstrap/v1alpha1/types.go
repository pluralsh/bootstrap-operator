package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&Bootstrap{}, &BootstrapList{})
}

type BootstrapSpec struct {
	ClusterName string               `json:"clusterName"`
	Components  ClusterAPIComponents `json:"clusterAPIComponents"`
	CloudSpec   CloudSpec            `json:"cloudSpec"`
}

type ClusterAPIComponents struct {
	Operator     ClusterAPIOperator     `json:"operator"`
	Core         ClusterAPICore         `json:"core"`
	ControlPlane ClusterAPIControlPlane `json:"controlPlane"`
	Bootstrap    ClusterAPIBootstrap    `json:"bootstrap"`
}

type ClusterAPIOperator struct {
	ManagerImage       string `json:"managerImage"`
	KubeRBACProxyImage string `json:"kubeRBACProxyImage"`
}

type ClusterAPICore struct {
	ClusterAPIComponentSpec `json:",inline"`
}

type ClusterAPIControlPlane struct {
	ClusterAPIComponentSpec `json:",inline"`
}

type ClusterAPIBootstrap struct {
	ClusterAPIComponentSpec `json:",inline"`
}

type ClusterAPIComponentSpec struct {
	Version        string `json:"version,omitempty"`
	FetchConfigURL string `json:"fetchConfigUrl,omitempty"`
}

type CloudSpec struct {
	AWS *AWSCloudSpec `json:"aws,omitempty"`
}

type AWSCloudSpec struct {
	Region                  string `json:"region"`
	ClusterAPIComponentSpec `json:",inline"`

	AccessKeyIDRef     corev1.SecretKeySelector `json:"accessKeyIdRef"`
	SecretAccessKeyRef corev1.SecretKeySelector `json:"secretAccessKeyRef"`
}

type BootstrapStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	CapiOperatorStatus *Status `json:"capiOperatorStatus,omitempty"`
	CapiBootstrap      *Status `json:"capiBootstrapStatus,omitempty"`
	CapiCore           *Status `json:"capiCoreStatus,omitempty"`
}

type Status struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Phase is a description of the current status, summarizing the various conditions.
	// This field is for informational purpose only and no logic
	// should be tied to the phase.
	// +optional
	Phase ComponentPhase `json:"phase,omitempty"`
}

// +kubebuilder:validation:Enum=Creating;Running;Failed

type ComponentPhase string

// These are the valid phases of a project.
const (
	Creating ComponentPhase = "Creating"
	Running  ComponentPhase = "Running"
	Failed   ComponentPhase = "Failed"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Application ready status"
type Bootstrap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BootstrapSpec   `json:"spec,omitempty"`
	Status BootstrapStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BootstrapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bootstrap `json:"items"`
}
