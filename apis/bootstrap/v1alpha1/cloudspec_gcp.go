package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type GCPCloudSpec struct {
	ClusterAPIComponentSpec `json:",inline"`
	CredentialsRef          corev1.SecretKeySelector `json:"credentialsRef"`
	// TODO
	Region string `json:"region"`
	// TODO
	Project string `json:"project"`

	// TODO
	Cluster *GCPManagedClusterSpec `json:"cluster"`

	// TODO
	// +optional
	MachinePool *GCPMachinePoolSpec `json:"machinePool"`

	// TODO
	// +optional
	ControlPlane *GCPManagedControlPlaneSpec `json:"controlPlane"`
}

type GCPManagedClusterSpec struct {
	// GCPNetworkSpec encapsulates all things related to the GCP network.
	// +optional
	Network GCPNetworkSpec `json:"network"`
}

type GCPNetworkSpec struct {
	// Name is the name of the network to be used.
	Name *string `json:"name,omitempty"`

	// AutoCreateSubnetworks: When set to true, the VPC network is created
	// in "auto" mode. When set to false, the VPC network is created in
	// "custom" mode.
	//
	// An auto mode VPC network starts with one subnet per region. Each
	// subnet has a predetermined range as described in Auto mode VPC
	// network IP ranges.
	//
	// Defaults to true.
	// +optional
	AutoCreateSubnetworks *bool `json:"autoCreateSubnetworks,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets GCPSubnets `json:"subnets,omitempty"`
}

// GCPSubnets is a slice of Subnet.
type GCPSubnets []GCPSubnetSpec

// GCPSubnetSpec configures an GCP Subnet.
type GCPSubnetSpec struct {
	// Name defines a unique identifier to reference this resource.
	Name string `json:"name,omitempty"`

	// CidrBlock is the range of internal addresses that are owned by this
	// subnetwork. Provide this property when you create the subnetwork. For
	// example, 10.0.0.0/8 or 192.168.0.0/16. Ranges must be unique and
	// non-overlapping within a network. Only IPv4 is supported. This field
	// can be set only at resource creation time.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// SecondaryCidrBlocks defines secondary CIDR ranges,
	// from which secondary IP ranges of a VM may be allocated
	// +optional
	SecondaryCidrBlocks map[string]string `json:"secondaryCidrBlocks,omitempty"`
}

type GCPMachinePoolSpec struct {
	// Replicas TODO
	Replicas int32 `json:"replicas"`
	// Scaling specifies scaling for the node pool
	// +optional
	Scaling *GCPNodePoolAutoscaling `json:"scaling,omitempty"`
}

// GCPNodePoolAutoscaling specifies scaling options.
type GCPNodePoolAutoscaling struct {
	MinCount *int32 `json:"minCount,omitempty"`
	MaxCount *int32 `json:"maxCount,omitempty"`
}

type GCPManagedControlPlaneSpec struct {
	// EnableAutopilot indicates whether to enable autopilot for this GKE cluster.
	// +optional
	EnableAutopilot bool `json:"enableAutopilot"`
	// EnableWorkloadIdentity allows enabling workload identity during cluster creation when
	// EnableAutopilot is disabled.
	EnableWorkloadIdentity bool `json:"enableWorkloadIdentity"`
	// ReleaseChannel represents the release channel of the GKE cluster.
	// "No channel" is used if ReleaseChannel is not set.
	// +optional
	ReleaseChannel *GCPReleaseChannel `json:"releaseChannel,omitempty"`
}

// GCPReleaseChannel is the release channel of the GKE cluster.
// +kubebuilder:validation:Enum=rapid;regular;stable
type GCPReleaseChannel string

const (
	// Rapid release channel.
	Rapid GCPReleaseChannel = "rapid"
	// Regular release channel.
	Regular GCPReleaseChannel = "regular"
	// Stable release channel.
	Stable GCPReleaseChannel = "stable"
)
