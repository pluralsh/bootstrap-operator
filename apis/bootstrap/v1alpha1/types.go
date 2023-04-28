package v1alpha1

import (
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
	SkipClusterCreation bool                     `json:"skipClusterCreation"`
	ClusterNetwork      ClusterNetwork           `json:"clusterNetwork"`
	ClusterAPI          ClusterAPI               `json:"clusterAPI"`
	CloudSpec           CloudSpec                `json:"cloudSpec"`
	GitHubSecretRef     corev1.SecretKeySelector `json:"gitHubSecretRef"`
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

type AWSCloudSpec struct {
	Region                  string `json:"region"`
	MachinePoolReplicas     int32  `json:"machinePoolReplicas"`
	InstanceType            string `json:"instanceType"`
	ClusterAPIComponentSpec `json:",inline"`
	// ServiceAccounts specifies service accounts
	// +optional
	ServiceAccounts []ClusterIAMServiceAccount `json:"serviceAccounts,omitempty"`
	// +optional
	Addons             []Addon                  `json:"addons,omitempty"`
	AccessKeyIDRef     corev1.SecretKeySelector `json:"accessKeyIdRef"`
	SecretAccessKeyRef corev1.SecretKeySelector `json:"secretAccessKeyRef"`
	SessionTokenRef    corev1.SecretKeySelector `json:"sessionTokenRef"`
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

type ClusterIAMServiceAccount struct {
	api.ClusterIAMMeta `json:"metadata,omitempty"`

	// list of ARNs of the IAM policies to attach
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
	// +optional
	WellKnownPolicies WellKnownPolicies `json:"wellKnownPolicies,omitempty"`

	// AttachPolicy holds a policy document to attach to this service account
	// +optional
	AttachPolicy string `json:"attachPolicy,omitempty"`

	// ARN of the role to attach to the service account
	// +optional
	AttachRoleARN string `json:"attachRoleARN,omitempty"`

	// ARN of the permissions boundary to associate with the service account
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`

	// Specific role name instead of the Cloudformation-generated role name
	// +optional
	RoleName string `json:"roleName,omitempty"`

	// Specify if only the IAM Service Account role should be created without creating/annotating the service account
	// +optional
	RoleOnly bool `json:"roleOnly,omitempty"`

	// AWS tags for the service account
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// WellKnownPolicies for attaching common IAM policies
type WellKnownPolicies struct {
	// ImageBuilder allows for full ECR (Elastic Container Registry) access.
	// +optional
	ImageBuilder bool `json:"imageBuilder,omitempty"`
	// AutoScaler adds policies for cluster-autoscaler. See [autoscaler AWS
	// docs](https://docs.aws.amazon.com/eks/latest/userguide/cluster-autoscaler.html).
	// +optional
	AutoScaler bool `json:"autoScaler,omitempty"`
	// AWSLoadBalancerController adds policies for using the
	// aws-load-balancer-controller. See [Load Balancer
	// docs](https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html).
	// +optional
	AWSLoadBalancerController bool `json:"awsLoadBalancerController,omitempty"`
	// ExternalDNS adds external-dns policies for Amazon Route 53.
	// See [external-dns
	// docs](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md).
	// +optional
	ExternalDNS bool `json:"externalDNS,omitempty"`
	// CertManager adds cert-manager policies. See [cert-manager
	// docs](https://cert-manager.io/docs/configuration/acme/dns01/route53).
	// +optional
	CertManager bool `json:"certManager,omitempty"`
	// EBSCSIController adds policies for using the
	// ebs-csi-controller. See [aws-ebs-csi-driver
	// docs](https://github.com/kubernetes-sigs/aws-ebs-csi-driver#set-up-driver-permission).
	// +optional
	EBSCSIController bool `json:"ebsCSIController,omitempty"`
	// EFSCSIController adds policies for using the
	// efs-csi-controller. See [aws-efs-csi-driver
	// docs](https://aws.amazon.com/blogs/containers/introducing-efs-csi-dynamic-provisioning).
	// +optional
	EFSCSIController bool `json:"efsCSIController,omitempty"`
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
