//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/errors"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIEndpoint) DeepCopyInto(out *APIEndpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIEndpoint.
func (in *APIEndpoint) DeepCopy() *APIEndpoint {
	if in == nil {
		return nil
	}
	out := new(APIEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSCloudSpec) DeepCopyInto(out *AWSCloudSpec) {
	*out = *in
	if in.MachinePools != nil {
		in, out := &in.MachinePools, &out.MachinePools
		*out = make([]AWSMachinePool, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
	in.IAMServiceAccount.DeepCopyInto(&out.IAMServiceAccount)
	if in.Addons != nil {
		in, out := &in.Addons, &out.Addons
		*out = make([]Addon, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.AccessKeyIDRef.DeepCopyInto(&out.AccessKeyIDRef)
	in.SecretAccessKeyRef.DeepCopyInto(&out.SecretAccessKeyRef)
	in.SessionTokenRef.DeepCopyInto(&out.SessionTokenRef)
	in.AWSAccountIDRef.DeepCopyInto(&out.AWSAccountIDRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSCloudSpec.
func (in *AWSCloudSpec) DeepCopy() *AWSCloudSpec {
	if in == nil {
		return nil
	}
	out := new(AWSCloudSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSMachinePool) DeepCopyInto(out *AWSMachinePool) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.AvailabilityZones != nil {
		in, out := &in.AvailabilityZones, &out.AvailabilityZones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SubnetIDs != nil {
		in, out := &in.SubnetIDs, &out.SubnetIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalTags != nil {
		in, out := &in.AdditionalTags, &out.AdditionalTags
		*out = make(Tags, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.RoleAdditionalPolicies != nil {
		in, out := &in.RoleAdditionalPolicies, &out.RoleAdditionalPolicies
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AMIVersion != nil {
		in, out := &in.AMIVersion, &out.AMIVersion
		*out = new(string)
		**out = **in
	}
	if in.AMIType != nil {
		in, out := &in.AMIType, &out.AMIType
		*out = new(ManagedMachineAMIType)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Taints != nil {
		in, out := &in.Taints, &out.Taints
		*out = make(Taints, len(*in))
		copy(*out, *in)
	}
	if in.DiskSize != nil {
		in, out := &in.DiskSize, &out.DiskSize
		*out = new(int32)
		**out = **in
	}
	if in.InstanceType != nil {
		in, out := &in.InstanceType, &out.InstanceType
		*out = new(string)
		**out = **in
	}
	if in.Scaling != nil {
		in, out := &in.Scaling, &out.Scaling
		*out = new(ManagedMachinePoolScaling)
		(*in).DeepCopyInto(*out)
	}
	if in.RemoteAccess != nil {
		in, out := &in.RemoteAccess, &out.RemoteAccess
		*out = new(ManagedRemoteAccess)
		(*in).DeepCopyInto(*out)
	}
	if in.ProviderIDList != nil {
		in, out := &in.ProviderIDList, &out.ProviderIDList
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CapacityType != nil {
		in, out := &in.CapacityType, &out.CapacityType
		*out = new(ManagedMachinePoolCapacityType)
		**out = **in
	}
	if in.UpdateConfig != nil {
		in, out := &in.UpdateConfig, &out.UpdateConfig
		*out = new(UpdateConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSMachinePool.
func (in *AWSMachinePool) DeepCopy() *AWSMachinePool {
	if in == nil {
		return nil
	}
	out := new(AWSMachinePool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Addon) DeepCopyInto(out *Addon) {
	*out = *in
	if in.ServiceAccountRoleArn != nil {
		in, out := &in.ServiceAccountRoleArn, &out.ServiceAccountRoleArn
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Addon.
func (in *Addon) DeepCopy() *Addon {
	if in == nil {
		return nil
	}
	out := new(Addon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllowedNamespaces) DeepCopyInto(out *AllowedNamespaces) {
	*out = *in
	if in.NamespaceList != nil {
		in, out := &in.NamespaceList, &out.NamespaceList
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllowedNamespaces.
func (in *AllowedNamespaces) DeepCopy() *AllowedNamespaces {
	if in == nil {
		return nil
	}
	out := new(AllowedNamespaces)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureCloudSpec) DeepCopyInto(out *AzureCloudSpec) {
	*out = *in
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
	in.ClusterIdentity.DeepCopyInto(&out.ClusterIdentity)
	if in.ManagedCluster != nil {
		in, out := &in.ManagedCluster, &out.ManagedCluster
		*out = new(AzureManagedClusterSpec)
		**out = **in
	}
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(AzureManagedControlPlaneSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.MachinePools != nil {
		in, out := &in.MachinePools, &out.MachinePools
		*out = make([]*AzureMachinePool, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(AzureMachinePool)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureCloudSpec.
func (in *AzureCloudSpec) DeepCopy() *AzureCloudSpec {
	if in == nil {
		return nil
	}
	out := new(AzureCloudSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureClusterIdentity) DeepCopyInto(out *AzureClusterIdentity) {
	*out = *in
	out.ClientSecret = in.ClientSecret
	if in.AllowedNamespaces != nil {
		in, out := &in.AllowedNamespaces, &out.AllowedNamespaces
		*out = new(AllowedNamespaces)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureClusterIdentity.
func (in *AzureClusterIdentity) DeepCopy() *AzureClusterIdentity {
	if in == nil {
		return nil
	}
	out := new(AzureClusterIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureMachinePool) DeepCopyInto(out *AzureMachinePool) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.ScaleSetPriority != nil {
		in, out := &in.ScaleSetPriority, &out.ScaleSetPriority
		*out = new(string)
		**out = **in
	}
	if in.Scaling != nil {
		in, out := &in.Scaling, &out.Scaling
		*out = new(ManagedMachinePoolScaling)
		(*in).DeepCopyInto(*out)
	}
	if in.AvailabilityZones != nil {
		in, out := &in.AvailabilityZones, &out.AvailabilityZones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.OsDiskType != nil {
		in, out := &in.OsDiskType, &out.OsDiskType
		*out = new(string)
		**out = **in
	}
	if in.OSDiskSizeGB != nil {
		in, out := &in.OSDiskSizeGB, &out.OSDiskSizeGB
		*out = new(int32)
		**out = **in
	}
	if in.MaxPods != nil {
		in, out := &in.MaxPods, &out.MaxPods
		*out = new(int32)
		**out = **in
	}
	if in.NodeLabels != nil {
		in, out := &in.NodeLabels, &out.NodeLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Taints != nil {
		in, out := &in.Taints, &out.Taints
		*out = make(Taints, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalTags != nil {
		in, out := &in.AdditionalTags, &out.AdditionalTags
		*out = make(Tags, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureMachinePool.
func (in *AzureMachinePool) DeepCopy() *AzureMachinePool {
	if in == nil {
		return nil
	}
	out := new(AzureMachinePool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureManagedClusterSpec) DeepCopyInto(out *AzureManagedClusterSpec) {
	*out = *in
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureManagedClusterSpec.
func (in *AzureManagedClusterSpec) DeepCopy() *AzureManagedClusterSpec {
	if in == nil {
		return nil
	}
	out := new(AzureManagedClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureManagedControlPlaneSpec) DeepCopyInto(out *AzureManagedControlPlaneSpec) {
	*out = *in
	if in.IdentityRef != nil {
		in, out := &in.IdentityRef, &out.IdentityRef
		*out = new(corev1.ObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureManagedControlPlaneSpec.
func (in *AzureManagedControlPlaneSpec) DeepCopy() *AzureManagedControlPlaneSpec {
	if in == nil {
		return nil
	}
	out := new(AzureManagedControlPlaneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Bootstrap) DeepCopyInto(out *Bootstrap) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Bootstrap.
func (in *Bootstrap) DeepCopy() *Bootstrap {
	if in == nil {
		return nil
	}
	out := new(Bootstrap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Bootstrap) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BootstrapList) DeepCopyInto(out *BootstrapList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Bootstrap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BootstrapList.
func (in *BootstrapList) DeepCopy() *BootstrapList {
	if in == nil {
		return nil
	}
	out := new(BootstrapList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BootstrapList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BootstrapSpec) DeepCopyInto(out *BootstrapSpec) {
	*out = *in
	in.ClusterNetwork.DeepCopyInto(&out.ClusterNetwork)
	in.ClusterAPI.DeepCopyInto(&out.ClusterAPI)
	in.CloudSpec.DeepCopyInto(&out.CloudSpec)
	if in.GitHubSecretRef != nil {
		in, out := &in.GitHubSecretRef, &out.GitHubSecretRef
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BootstrapSpec.
func (in *BootstrapSpec) DeepCopy() *BootstrapSpec {
	if in == nil {
		return nil
	}
	out := new(BootstrapSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BootstrapStatus) DeepCopyInto(out *BootstrapStatus) {
	*out = *in
	out.Status = in.Status
	if in.CapiOperatorStatus != nil {
		in, out := &in.CapiOperatorStatus, &out.CapiOperatorStatus
		*out = new(Status)
		**out = **in
	}
	if in.CapiOperatorComponentsStatus != nil {
		in, out := &in.CapiOperatorComponentsStatus, &out.CapiOperatorComponentsStatus
		*out = new(Status)
		**out = **in
	}
	if in.CapiClusterStatus != nil {
		in, out := &in.CapiClusterStatus, &out.CapiClusterStatus
		*out = new(ClusterStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.ProviderStatus != nil {
		in, out := &in.ProviderStatus, &out.ProviderStatus
		*out = new(Status)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BootstrapStatus.
func (in *BootstrapStatus) DeepCopy() *BootstrapStatus {
	if in == nil {
		return nil
	}
	out := new(BootstrapStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudSpec) DeepCopyInto(out *CloudSpec) {
	*out = *in
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(AWSCloudSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureCloudSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.GCP != nil {
		in, out := &in.GCP, &out.GCP
		*out = new(GCPCloudSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudSpec.
func (in *CloudSpec) DeepCopy() *CloudSpec {
	if in == nil {
		return nil
	}
	out := new(CloudSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPI) DeepCopyInto(out *ClusterAPI) {
	*out = *in
	in.Components.DeepCopyInto(&out.Components)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPI.
func (in *ClusterAPI) DeepCopy() *ClusterAPI {
	if in == nil {
		return nil
	}
	out := new(ClusterAPI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPIBootstrap) DeepCopyInto(out *ClusterAPIBootstrap) {
	*out = *in
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPIBootstrap.
func (in *ClusterAPIBootstrap) DeepCopy() *ClusterAPIBootstrap {
	if in == nil {
		return nil
	}
	out := new(ClusterAPIBootstrap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPIComponentSpec) DeepCopyInto(out *ClusterAPIComponentSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPIComponentSpec.
func (in *ClusterAPIComponentSpec) DeepCopy() *ClusterAPIComponentSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterAPIComponentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPIComponents) DeepCopyInto(out *ClusterAPIComponents) {
	*out = *in
	out.Operator = in.Operator
	if in.Core != nil {
		in, out := &in.Core, &out.Core
		*out = new(ClusterAPICore)
		**out = **in
	}
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(ClusterAPIControlPlane)
		**out = **in
	}
	if in.Bootstrap != nil {
		in, out := &in.Bootstrap, &out.Bootstrap
		*out = new(ClusterAPIBootstrap)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPIComponents.
func (in *ClusterAPIComponents) DeepCopy() *ClusterAPIComponents {
	if in == nil {
		return nil
	}
	out := new(ClusterAPIComponents)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPIControlPlane) DeepCopyInto(out *ClusterAPIControlPlane) {
	*out = *in
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPIControlPlane.
func (in *ClusterAPIControlPlane) DeepCopy() *ClusterAPIControlPlane {
	if in == nil {
		return nil
	}
	out := new(ClusterAPIControlPlane)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPICore) DeepCopyInto(out *ClusterAPICore) {
	*out = *in
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPICore.
func (in *ClusterAPICore) DeepCopy() *ClusterAPICore {
	if in == nil {
		return nil
	}
	out := new(ClusterAPICore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAPIOperator) DeepCopyInto(out *ClusterAPIOperator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAPIOperator.
func (in *ClusterAPIOperator) DeepCopy() *ClusterAPIOperator {
	if in == nil {
		return nil
	}
	out := new(ClusterAPIOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterIAMMeta) DeepCopyInto(out *ClusterIAMMeta) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterIAMMeta.
func (in *ClusterIAMMeta) DeepCopy() *ClusterIAMMeta {
	if in == nil {
		return nil
	}
	out := new(ClusterIAMMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterIAMServiceAccount) DeepCopyInto(out *ClusterIAMServiceAccount) {
	*out = *in
	in.ClusterIAMMeta.DeepCopyInto(&out.ClusterIAMMeta)
	if in.AttachPolicyARNs != nil {
		in, out := &in.AttachPolicyARNs, &out.AttachPolicyARNs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.WellKnownPolicies = in.WellKnownPolicies
	if in.AttachPolicy != nil {
		in, out := &in.AttachPolicy, &out.AttachPolicy
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterIAMServiceAccount.
func (in *ClusterIAMServiceAccount) DeepCopy() *ClusterIAMServiceAccount {
	if in == nil {
		return nil
	}
	out := new(ClusterIAMServiceAccount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterNetwork) DeepCopyInto(out *ClusterNetwork) {
	*out = *in
	if in.APIServerPort != nil {
		in, out := &in.APIServerPort, &out.APIServerPort
		*out = new(int32)
		**out = **in
	}
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = new(NetworkRanges)
		(*in).DeepCopyInto(*out)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = new(NetworkRanges)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterNetwork.
func (in *ClusterNetwork) DeepCopy() *ClusterNetwork {
	if in == nil {
		return nil
	}
	out := new(ClusterNetwork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatus) DeepCopyInto(out *ClusterStatus) {
	*out = *in
	out.Status = in.Status
	if in.FailureReason != nil {
		in, out := &in.FailureReason, &out.FailureReason
		*out = new(errors.ClusterStatusError)
		**out = **in
	}
	if in.FailureMessage != nil {
		in, out := &in.FailureMessage, &out.FailureMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatus.
func (in *ClusterStatus) DeepCopy() *ClusterStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPCloudSpec) DeepCopyInto(out *GCPCloudSpec) {
	*out = *in
	out.ClusterAPIComponentSpec = in.ClusterAPIComponentSpec
	in.CredentialsRef.DeepCopyInto(&out.CredentialsRef)
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(GCPManagedClusterSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.MachinePool != nil {
		in, out := &in.MachinePool, &out.MachinePool
		*out = new(GCPMachinePoolSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(GCPManagedControlPlaneSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPCloudSpec.
func (in *GCPCloudSpec) DeepCopy() *GCPCloudSpec {
	if in == nil {
		return nil
	}
	out := new(GCPCloudSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPMachinePoolSpec) DeepCopyInto(out *GCPMachinePoolSpec) {
	*out = *in
	if in.Scaling != nil {
		in, out := &in.Scaling, &out.Scaling
		*out = new(GCPNodePoolAutoscaling)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPMachinePoolSpec.
func (in *GCPMachinePoolSpec) DeepCopy() *GCPMachinePoolSpec {
	if in == nil {
		return nil
	}
	out := new(GCPMachinePoolSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPManagedClusterSpec) DeepCopyInto(out *GCPManagedClusterSpec) {
	*out = *in
	in.Network.DeepCopyInto(&out.Network)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPManagedClusterSpec.
func (in *GCPManagedClusterSpec) DeepCopy() *GCPManagedClusterSpec {
	if in == nil {
		return nil
	}
	out := new(GCPManagedClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPManagedControlPlaneSpec) DeepCopyInto(out *GCPManagedControlPlaneSpec) {
	*out = *in
	if in.ReleaseChannel != nil {
		in, out := &in.ReleaseChannel, &out.ReleaseChannel
		*out = new(GCPReleaseChannel)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPManagedControlPlaneSpec.
func (in *GCPManagedControlPlaneSpec) DeepCopy() *GCPManagedControlPlaneSpec {
	if in == nil {
		return nil
	}
	out := new(GCPManagedControlPlaneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPNetworkSpec) DeepCopyInto(out *GCPNetworkSpec) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.AutoCreateSubnetworks != nil {
		in, out := &in.AutoCreateSubnetworks, &out.AutoCreateSubnetworks
		*out = new(bool)
		**out = **in
	}
	if in.Subnets != nil {
		in, out := &in.Subnets, &out.Subnets
		*out = make(GCPSubnets, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPNetworkSpec.
func (in *GCPNetworkSpec) DeepCopy() *GCPNetworkSpec {
	if in == nil {
		return nil
	}
	out := new(GCPNetworkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPNodePoolAutoscaling) DeepCopyInto(out *GCPNodePoolAutoscaling) {
	*out = *in
	if in.MinCount != nil {
		in, out := &in.MinCount, &out.MinCount
		*out = new(int32)
		**out = **in
	}
	if in.MaxCount != nil {
		in, out := &in.MaxCount, &out.MaxCount
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPNodePoolAutoscaling.
func (in *GCPNodePoolAutoscaling) DeepCopy() *GCPNodePoolAutoscaling {
	if in == nil {
		return nil
	}
	out := new(GCPNodePoolAutoscaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPSubnetSpec) DeepCopyInto(out *GCPSubnetSpec) {
	*out = *in
	if in.SecondaryCidrBlocks != nil {
		in, out := &in.SecondaryCidrBlocks, &out.SecondaryCidrBlocks
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPSubnetSpec.
func (in *GCPSubnetSpec) DeepCopy() *GCPSubnetSpec {
	if in == nil {
		return nil
	}
	out := new(GCPSubnetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in GCPSubnets) DeepCopyInto(out *GCPSubnets) {
	{
		in := &in
		*out = make(GCPSubnets, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPSubnets.
func (in GCPSubnets) DeepCopy() GCPSubnets {
	if in == nil {
		return nil
	}
	out := new(GCPSubnets)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IAMServiceAccountSpec) DeepCopyInto(out *IAMServiceAccountSpec) {
	*out = *in
	if in.RoleNamePrefix != nil {
		in, out := &in.RoleNamePrefix, &out.RoleNamePrefix
		*out = new(string)
		**out = **in
	}
	if in.ServiceAccounts != nil {
		in, out := &in.ServiceAccounts, &out.ServiceAccounts
		*out = make([]ClusterIAMServiceAccount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IAMServiceAccountSpec.
func (in *IAMServiceAccountSpec) DeepCopy() *IAMServiceAccountSpec {
	if in == nil {
		return nil
	}
	out := new(IAMServiceAccountSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedMachinePoolScaling) DeepCopyInto(out *ManagedMachinePoolScaling) {
	*out = *in
	if in.MinSize != nil {
		in, out := &in.MinSize, &out.MinSize
		*out = new(int32)
		**out = **in
	}
	if in.MaxSize != nil {
		in, out := &in.MaxSize, &out.MaxSize
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedMachinePoolScaling.
func (in *ManagedMachinePoolScaling) DeepCopy() *ManagedMachinePoolScaling {
	if in == nil {
		return nil
	}
	out := new(ManagedMachinePoolScaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedRemoteAccess) DeepCopyInto(out *ManagedRemoteAccess) {
	*out = *in
	if in.SSHKeyName != nil {
		in, out := &in.SSHKeyName, &out.SSHKeyName
		*out = new(string)
		**out = **in
	}
	if in.SourceSecurityGroups != nil {
		in, out := &in.SourceSecurityGroups, &out.SourceSecurityGroups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedRemoteAccess.
func (in *ManagedRemoteAccess) DeepCopy() *ManagedRemoteAccess {
	if in == nil {
		return nil
	}
	out := new(ManagedRemoteAccess)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkRanges) DeepCopyInto(out *NetworkRanges) {
	*out = *in
	if in.CIDRBlocks != nil {
		in, out := &in.CIDRBlocks, &out.CIDRBlocks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkRanges.
func (in *NetworkRanges) DeepCopy() *NetworkRanges {
	if in == nil {
		return nil
	}
	out := new(NetworkRanges)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Tags) DeepCopyInto(out *Tags) {
	{
		in := &in
		*out = make(Tags, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tags.
func (in Tags) DeepCopy() Tags {
	if in == nil {
		return nil
	}
	out := new(Tags)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Taint) DeepCopyInto(out *Taint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Taint.
func (in *Taint) DeepCopy() *Taint {
	if in == nil {
		return nil
	}
	out := new(Taint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Taints) DeepCopyInto(out *Taints) {
	{
		in := &in
		*out = make(Taints, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Taints.
func (in Taints) DeepCopy() Taints {
	if in == nil {
		return nil
	}
	out := new(Taints)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpdateConfig) DeepCopyInto(out *UpdateConfig) {
	*out = *in
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(int)
		**out = **in
	}
	if in.MaxUnavailablePercentage != nil {
		in, out := &in.MaxUnavailablePercentage, &out.MaxUnavailablePercentage
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateConfig.
func (in *UpdateConfig) DeepCopy() *UpdateConfig {
	if in == nil {
		return nil
	}
	out := new(UpdateConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WellKnownPolicies) DeepCopyInto(out *WellKnownPolicies) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WellKnownPolicies.
func (in *WellKnownPolicies) DeepCopy() *WellKnownPolicies {
	if in == nil {
		return nil
	}
	out := new(WellKnownPolicies)
	in.DeepCopyInto(out)
	return out
}
