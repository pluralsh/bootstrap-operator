package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AzureCloudSpec struct {
	ClusterAPIComponentSpec `json:",inline"`
	ClusterIdentity         AzureClusterIdentity          `json:"clusterIdentity"`
	ManagedCluster          *AzureManagedClusterSpec      `json:"managedCluster"`
	ControlPlane            *AzureManagedControlPlaneSpec `json:"controlPlane"`
	MachinePools            []*AzureMachinePool           `json:"machinePools"`
}

// AzureClusterIdentity defines the parameters that are used to create an AzureIdentity.
type AzureClusterIdentity struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Type is the type of Azure Identity used.
	// ServicePrincipal, ServicePrincipalCertificate, UserAssignedMSI or ManualServicePrincipal.
	Type IdentityType `json:"type"`
	// ResourceID is the Azure resource ID for the User Assigned MSI resource.
	// Only applicable when type is UserAssignedMSI.
	// +optional
	ResourceID string `json:"resourceID,omitempty"`
	// ClientID is the service principal client ID.
	// Both User Assigned MSI and SP can use this field.
	ClientID string `json:"clientID"`
	// ClientSecret is a secret reference which should contain either a Service Principal password or certificate secret.
	// +optional
	ClientSecret corev1.SecretReference `json:"clientSecret,omitempty"`
	// TenantID is the service principal primary tenant id.
	TenantID string `json:"tenantID"`
	// AllowedNamespaces is used to identify the namespaces the clusters are allowed to use the identity from.
	// Namespaces can be selected either using an array of namespaces or with label selector.
	// An empty allowedNamespaces object indicates that AzureClusters can use this identity from any namespace.
	// If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
	// A namespace should be either in the NamespaceList or match with Selector to use the identity.
	//
	// +optional
	// +nullable
	AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces"`
}

// IdentityType represents different types of identities.
// +kubebuilder:validation:Enum=ServicePrincipal;UserAssignedMSI;ManualServicePrincipal;ServicePrincipalCertificate
type IdentityType string

// AllowedNamespaces defines the namespaces the clusters are allowed to use the identity from
// NamespaceList takes precedence over the Selector.
type AllowedNamespaces struct {
	// A nil or empty list indicates that AzureCluster cannot use the identity from any namespace.
	//
	// +optional
	// +nullable
	NamespaceList []string `json:"list"`
	// Selector is a selector of namespaces that AzureCluster can
	// use this Identity from. This is a standard Kubernetes LabelSelector,
	// a label query over a set of resources. The result of matchLabels and
	// matchExpressions are ANDed.
	//
	// A nil or empty selector indicates that AzureCluster cannot use this
	// AzureClusterIdentity from any namespace.
	// +optional
	Selector *metav1.LabelSelector `json:"selector"`
}

// AzureManagedClusterSpec defines the desired state of AzureManagedCluster.
type AzureManagedClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint"`
}

// AzureManagedControlPlaneSpec defines the desired state of AzureManagedControlPlane.
type AzureManagedControlPlaneSpec struct {
	// Version defines the desired Kubernetes version.
	// +kubebuilder:validation:MinLength:=2
	Version string `json:"version"`

	// ResourceGroupName is the name of the Azure resource group for this AKS Cluster.
	ResourceGroupName string `json:"resourceGroupName"`

	// SubscriptionID is the GUID of the Azure subscription to hold this cluster.
	SubscriptionID string `json:"subscriptionID,omitempty"`

	// Location is a string matching one of the canonical Azure region names. Examples: "westus2", "eastus".
	Location string `json:"location"`

	// SSHPublicKey is a string literal containing an ssh public key base64 encoded.
	SSHPublicKey string `json:"sshPublicKey"`

	// IdentityRef is a reference to a AzureClusterIdentity to be used when reconciling this cluster
	IdentityRef *corev1.ObjectReference `json:"identityRef,omitempty"`
}

type AzureMachinePool struct {
	// Number of desired machines.
	Replicas *int32 `json:"replicas,omitempty"`

	// Name of the agent pool.
	Name string `json:"name,omitempty"`

	// Mode represents mode of an agent pool. Possible values include: System, User.
	// +kubebuilder:validation:Enum=System;User
	Mode string `json:"mode"`

	// SKU is the size of the VMs in the node pool.
	SKU string `json:"sku"`
}
