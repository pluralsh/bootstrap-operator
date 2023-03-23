package v1alpha1

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

func init() {
	SchemeBuilder.Register(&Bootstrap{}, &BootstrapList{})
}

type BootstrapSpec struct {
	ClusterName       string         `json:"clusterName"`
	KubernetesVersion string         `json:"kubernetesVersion"`
	ClusterNetwork    ClusterNetwork `json:"clusterNetwork"`
	ClusterAPI        ClusterAPI     `json:"clusterAPI"`
	CloudSpec         CloudSpec      `json:"cloudSpec"`
}

// ClusterNetwork specifies the different networking
// parameters for a cluster.
type ClusterNetwork struct {
	// APIServerPort specifies the port the API Server should bind to.
	// Defaults to 6443.
	// +optional
	APIServerPort *int32 `json:"apiServerPort,omitempty"`

	// The network ranges from which service VIPs are allocated.
	// +optional
	Services *NetworkRanges `json:"services,omitempty"`

	// The network ranges from which Pod networks are allocated.
	// +optional
	Pods *NetworkRanges `json:"pods,omitempty"`

	// Domain name for services.
	// +optional
	ServiceDomain string `json:"serviceDomain,omitempty"`
}

// NetworkRanges represents ranges of network addresses.
type NetworkRanges struct {
	CIDRBlocks []string `json:"cidrBlocks"`
}

func (n NetworkRanges) String() string {
	if len(n.CIDRBlocks) == 0 {
		return ""
	}
	return strings.Join(n.CIDRBlocks, ",")
}

type ClusterAPIComponents struct {
	Operator ClusterAPIOperator `json:"operator"`
}

type ClusterAPI struct {
	Components ClusterAPIComponents `json:"components"`
	Version    string               `json:"version"`
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
	MachinePoolReplicas     int32  `json:"machinePoolReplicas"`
	InstanceType            string `json:"instanceType"`
	ClusterAPIComponentSpec `json:",inline"`

	AccessKeyIDRef     corev1.SecretKeySelector `json:"accessKeyIdRef"`
	SecretAccessKeyRef corev1.SecretKeySelector `json:"secretAccessKeyRef"`
	SessionTokenRef    corev1.SecretKeySelector `json:"sessionTokenRef"`
}

type BootstrapStatus struct {
	Status `json:",inline"`

	CapiOperatorStatus           *Status        `json:"capiOperatorStatus,omitempty"`
	CapiOperatorComponentsStatus *Status        `json:"capiOperatorComponentsStatus,omitempty"`
	CapiClusterStatus            *ClusterStatus `json:"capiClusterStatus,omitempty"`
	ProviderStatus               *Status        `json:"providerStatus,omitempty"`
}

type ClusterStatus struct {
	Status `json:",inline"`
	// FailureReason indicates that there is a fatal problem reconciling the
	// state, and will be set to a token value suitable for
	// programmatic interpretation.
	// +optional
	FailureReason *capierrors.ClusterStatusError `json:"failureReason,omitempty"`

	// FailureMessage indicates that there is a fatal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// InfrastructureReady is the state of the infrastructure provider.
	// +optional
	InfrastructureReady bool `json:"infrastructureReady"`

	// ControlPlaneReady defines if the control plane is ready.
	// +optional
	ControlPlaneReady bool `json:"controlPlaneReady"`
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

// +kubebuilder:validation:Enum=Started;Creating;Running;Failed;Error

type ComponentPhase string

// These are the valid phases of a project.
const (
	Started  ComponentPhase = "Started"
	Creating ComponentPhase = "Creating"
	Running  ComponentPhase = "Running"
	Failed   ComponentPhase = "Failed"
	Error    ComponentPhase = "Error"
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
