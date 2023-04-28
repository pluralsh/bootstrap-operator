package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Tags map[string]string

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

var (
	// TaintEffectNoSchedule is a taint that indicates that a pod shouldn't be scheduled on a node
	// unless it can tolerate the taint.
	TaintEffectNoSchedule = TaintEffect("no-schedule")
	// TaintEffectNoExecute is a taint that indicates that a pod shouldn't be schedule on a node
	// unless it can tolerate it. And if its already running on the node it will be evicted.
	TaintEffectNoExecute = TaintEffect("no-execute")
	// TaintEffectPreferNoSchedule is a taint that indicates that there is a "preference" that pods shouldn't
	// be scheduled on a node unless it can tolerate the taint. the scheduler will try to avoid placing the pod
	// but it may still run on the node if there is no other option.
	TaintEffectPreferNoSchedule = TaintEffect("prefer-no-schedule")
)

// ManagedMachineAMIType specifies which AWS AMI to use for a managed MachinePool.
type ManagedMachineAMIType string

const (
	// Al2x86_64 is the default AMI type.
	Al2x86_64 ManagedMachineAMIType = "AL2_x86_64"
	// Al2x86_64GPU is the x86-64 GPU AMI type.
	Al2x86_64GPU ManagedMachineAMIType = "AL2_x86_64_GPU"
	// Al2Arm64 is the Arm AMI type.
	Al2Arm64 ManagedMachineAMIType = "AL2_ARM_64"
)

// ManagedMachinePoolCapacityType specifies the capacity type to be used for the managed MachinePool.
type ManagedMachinePoolCapacityType string

const (
	// ManagedMachinePoolCapacityTypeOnDemand is the default capacity type, to launch on-demand instances.
	ManagedMachinePoolCapacityTypeOnDemand ManagedMachinePoolCapacityType = "onDemand"
	// ManagedMachinePoolCapacityTypeSpot is the spot instance capacity type to launch spot instances.
	ManagedMachinePoolCapacityTypeSpot ManagedMachinePoolCapacityType = "spot"
)

type AWSCloudSpec struct {
	MachinePools            []AWSMachinePool `json:"machinePools"`
	Region                  string           `json:"region"`
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

type AWSMachinePool struct {
	Name string `json:"name"`
	// Number of desired machines. Defaults to 1.
	// This is a pointer to distinguish between explicit zero and not specified.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// EKSNodegroupName specifies the name of the nodegroup in AWS
	// corresponding to this MachinePool. If you don't specify a name
	// then a default name will be created based on the namespace and
	// name of the managed machine pool.
	// +optional
	EKSNodegroupName string `json:"eksNodegroupName,omitempty"`

	// AvailabilityZones is an array of availability zones instances can run in
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// SubnetIDs specifies which subnets are used for the
	// auto scaling group of this nodegroup
	// +optional
	SubnetIDs []string `json:"subnetIDs,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the node group role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies []string `json:"roleAdditionalPolicies,omitempty"`

	// RoleName specifies the name of IAM role for the node group.
	// If the role is pre-existing we will treat it as unmanaged
	// and not delete it on deletion. If the EKSEnableIAM feature
	// flag is true and no name is supplied then a role is created.
	// +optional
	RoleName string `json:"roleName,omitempty"`

	// AMIVersion defines the desired AMI release version. If no version number
	// is supplied then the latest version for the Kubernetes version
	// will be used
	// +kubebuilder:validation:MinLength:=2
	// +optional
	AMIVersion *string `json:"amiVersion,omitempty"`

	// AMIType defines the AMI type
	// +kubebuilder:validation:Enum:=AL2_x86_64;AL2_x86_64_GPU;AL2_ARM_64;CUSTOM
	// +kubebuilder:default:=AL2_x86_64
	// +optional
	AMIType *ManagedMachineAMIType `json:"amiType,omitempty"`

	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Taints specifies the taints to apply to the nodes of the machine pool
	// +optional
	Taints Taints `json:"taints,omitempty"`

	// DiskSize specifies the root disk size
	// +optional
	DiskSize *int32 `json:"diskSize,omitempty"`

	// InstanceType specifies the AWS instance type
	// +optional
	InstanceType *string `json:"instanceType,omitempty"`

	// Scaling specifies scaling for the ASG behind this pool
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`

	// RemoteAccess specifies how machines can be accessed remotely
	// +optional
	RemoteAccess *ManagedRemoteAccess `json:"remoteAccess,omitempty"`

	// ProviderIDList are the provider IDs of instances in the
	// autoscaling group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// CapacityType specifies the capacity type for the ASG behind this pool
	// +kubebuilder:validation:Enum:=onDemand;spot
	// +kubebuilder:default:=onDemand
	// +optional
	CapacityType *ManagedMachinePoolCapacityType `json:"capacityType,omitempty"`

	// UpdateConfig holds the optional config to control the behaviour of the update
	// to the nodegroup.
	// +optional
	UpdateConfig *UpdateConfig `json:"updateConfig,omitempty"`
}

// UpdateConfig is the configuration options for updating a nodegroup. Only one of MaxUnavailable
// and MaxUnavailablePercentage should be specified.
type UpdateConfig struct {
	// MaxUnavailable is the maximum number of nodes unavailable at once during a version update.
	// Nodes will be updated in parallel. The maximum number is 100.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailable *int `json:"maxUnavailable,omitempty"`

	// MaxUnavailablePercentage is the maximum percentage of nodes unavailable during a version update. This
	// percentage of nodes will be updated in parallel, up to 100 nodes at once.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailablePercentage *int `json:"maxUnavailablePercentage,omitempty"`
}

// ManagedRemoteAccess specifies remote access settings for EC2 instances.
type ManagedRemoteAccess struct {
	// SSHKeyName specifies which EC2 SSH key can be used to access machines.
	// If left empty, the key from the control plane is used.
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// SourceSecurityGroups specifies which security groups are allowed access
	SourceSecurityGroups []string `json:"sourceSecurityGroups,omitempty"`

	// Public specifies whether to open port 22 to the public internet
	Public bool `json:"public,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize *int32 `json:"minSize,omitempty"`
	MaxSize *int32 `json:"maxSize,omitempty"`
}

type ClusterIAMServiceAccount struct {
	ClusterIAMMeta `json:"metadata,omitempty"`

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

// ClusterIAMMeta holds information we can use to create ObjectMeta for service
// accounts
type ClusterIAMMeta struct {
	// +optional
	Name string `json:"name,omitempty"`

	// +optional
	Namespace string `json:"namespace,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AsObjectMeta gives us the k8s ObjectMeta needed to create the service account
func (iamMeta *ClusterIAMMeta) AsObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        iamMeta.Name,
		Namespace:   iamMeta.Namespace,
		Annotations: iamMeta.Annotations,
		Labels:      iamMeta.Labels,
	}
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
