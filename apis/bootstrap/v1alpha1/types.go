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
	ClusterName       string `json:"clusterName"`
	KubernetesVersion string `json:"kubernetesVersion"`
	// +kubebuilder:default:=false
	// +optional
	SkipClusterCreation bool `json:"skipClusterCreation"`
	// +kubebuilder:default:=false
	// +optional
	MoveCluster bool `json:"moveCluster"`
	// +kubebuilder:default:=false
	// +optional
	BootstrapMode bool `json:"bootstrapMode,omitempty"`
	// +kubebuilder:default:=false
	// +optional
	MigrateCluster bool           `json:"migrateCluster,omitempty"`
	ClusterNetwork ClusterNetwork `json:"clusterNetwork"`
	ClusterAPI     ClusterAPI     `json:"clusterAPI"`
	CloudSpec      CloudSpec      `json:"cloudSpec"`
	// +optional
	GitHubSecretRef *corev1.SecretKeySelector `json:"gitHubSecretRef,omitempty"`
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
	Operator     ClusterAPIOperator      `json:"operator"`
	Core         *ClusterAPICore         `json:"core,omitempty"`
	ControlPlane *ClusterAPIControlPlane `json:"controlPlane,omitempty"`
	Bootstrap    *ClusterAPIBootstrap    `json:"bootstrap,omitempty"`
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
	AWS   *AWSCloudSpec   `json:"aws,omitempty"`
	Azure *AzureCloudSpec `json:"azure,omitempty"`
	GCP   *GCPCloudSpec   `json:"gcp,omitempty"`
}

// Addon represents a EKS addon.
type Addon struct {
	// Name is the name of the addon
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Version is the version of the addon to use
	Version string `json:"version"`
	// ServiceAccountRoleArn is the ARN of an IAM role to bind to the addons service account
	// +optional
	ServiceAccountRoleArn *string `json:"serviceAccountRoleARN,omitempty"`
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

type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int32 `json:"port"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize *int32 `json:"minSize,omitempty"`
	MaxSize *int32 `json:"maxSize,omitempty"`
}

// Taint defines the specs for a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=no-schedule;no-execute;prefer-no-schedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// Value is the value of the taint
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// Equals is used to test if 2 taints are equal.
func (t *Taint) Equals(other *Taint) bool {
	if t == nil || other == nil {
		return t == other
	}

	return t.Effect == other.Effect &&
		t.Key == other.Key &&
		t.Value == other.Value
}

// Taints is an array of Taints.
type Taints []Taint

// Contains checks for existence of a matching taint.
func (t *Taints) Contains(taint *Taint) bool {
	for _, t := range *t {
		if t.Equals(taint) {
			return true
		}
	}

	return false
}
