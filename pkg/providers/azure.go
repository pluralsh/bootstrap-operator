package providers

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	azure "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	r "github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

type AzureProvider struct {
	Data           *resources.TemplateData
	GitHubToken    string
	version        string
	fetchConfigUrl string
}

const (
	azureSecretName = "azure-credentials"
)

func GetAzureProvider(data *resources.TemplateData) (*AzureProvider, error) {
	spec := data.Bootstrap.Spec
	azureSpec := spec.CloudSpec.Azure

	var gitHubToken string
	if spec.GitHubSecretRef != nil && spec.GitHubSecretRef.Name != "" && spec.GitHubSecretRef.Key != "" {
		var secret corev1.Secret
		if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{
			Namespace: data.Bootstrap.Namespace, Name: spec.GitHubSecretRef.Name}, &secret); err != nil {
			return nil, err
		}
		gitHubToken = strings.TrimSpace(string(secret.Data[spec.GitHubSecretRef.Key]))
	}

	data.Log.Named("Azure provider").Info("Create Azure provider")
	return &AzureProvider{
		Data:           data,
		version:        azureSpec.Version,
		fetchConfigUrl: azureSpec.FetchConfigURL,
		GitHubToken:    gitHubToken,
	}, nil
}

func (azure *AzureProvider) Name() string {
	return "azure"
}

func (azure *AzureProvider) Version() string {
	return azure.version
}

func (azure *AzureProvider) FetchConfigURL() string {
	return azure.fetchConfigUrl
}

func (azure *AzureProvider) createCredentialSecret() error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: azure.Data.Bootstrap.Namespace,
			Name:      azureSecretName,
		},
		Data: map[string][]byte{
			"EXP_MACHINE_POOL": []byte("true"),
			"GITHUB_TOKEN":     []byte(azure.GitHubToken),
		},
	}
	if err := azure.Data.Client.Create(azure.Data.Ctx, &secret); err != nil {
		return err
	}
	return nil
}

func (azure *AzureProvider) Init() (*ctrl.Result, error) {
	if azure.Data.Bootstrap.Status.ProviderStatus.Phase != bv1alpha1.Creating {
		if err := azure.updateProviderStatus(bv1alpha1.Creating, "init Azure provider", false); err != nil {
			return nil, err
		}
		if err := azure.createCredentialSecret(); err != nil {
			if err := azure.updateProviderStatus(bv1alpha1.Error, err.Error(), false); err != nil {
				return nil, err
			}
		}
		if err := azure.updateProviderStatus(bv1alpha1.Running, "Azure provider ready", true); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (azure *AzureProvider) updateProviderStatus(phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(azure.Data.Ctx, azure.Data.Client, azure.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.ProviderStatus == nil {
			c.Status.ProviderStatus = &bv1alpha1.Status{}
		}
		c.Status.ProviderStatus.Message = message
		c.Status.ProviderStatus.Phase = phase
		c.Status.ProviderStatus.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update bootstrap: %w", message, err)
	}

	return nil
}

func (azure *AzureProvider) Secret() string {
	return azureSecretName
}

func (azure *AzureProvider) CheckCluster() (*ctrl.Result, error) {
	var cluster clusterv1.Cluster
	if err := azure.Data.Client.Get(azure.Data.Ctx, ctrlruntimeclient.ObjectKey{
		Namespace: azure.Data.Bootstrap.Namespace,
		Name:      azure.Data.Bootstrap.Spec.ClusterName}, &cluster); err != nil {
		return nil, err
	}
	if err := azure.updateClusterStatus(cluster.Status); err != nil {
		return nil, err
	}
	for _, cond := range cluster.Status.Conditions {
		if cond.Type == clusterv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			return nil, azure.setStatusReady()
		}
	}
	azure.Data.Log.Named("Azure provider").Info("Waiting for Azure cluster to become ready")
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (azure *AzureProvider) setStatusReady() error {
	return helper.UpdateBootstrapStatus(azure.Data.Ctx, azure.Data.Client, azure.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiClusterStatus == nil {
			c.Status.CapiClusterStatus = &bv1alpha1.ClusterStatus{}
		}
		c.Status.CapiClusterStatus.Ready = true
	})
}

func (azure *AzureProvider) updateClusterStatus(status clusterv1.ClusterStatus) error {
	err := helper.UpdateBootstrapStatus(azure.Data.Ctx, azure.Data.Client, azure.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiClusterStatus == nil {
			c.Status.CapiClusterStatus = &bv1alpha1.ClusterStatus{}
		}
		c.Status.CapiClusterStatus.ControlPlaneReady = status.ControlPlaneReady
		c.Status.CapiClusterStatus.InfrastructureReady = status.InfrastructureReady
		c.Status.CapiClusterStatus.FailureMessage = status.FailureMessage
		c.Status.CapiClusterStatus.FailureReason = status.FailureReason

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on cluster status. Could not update bootstrap: %w", err)
	}

	return nil
}

func (azure *AzureProvider) ReconcileCluster() error {
	ctx := azure.Data.Ctx
	client := azure.Data.Client
	namespace := azure.Data.Bootstrap.Namespace

	clusterIdentityCreator := []r.NamedAzureClusterIdentityCreatorGetter{
		azureClusterIdentityCreator(azure.Data),
	}
	if err := r.ReconcileAzureClusterIdentitys(ctx, clusterIdentityCreator, namespace, client); err != nil {
		return err
	}

	clusterCreator := []r.NamedClusterCreatorGetter{
		azureClusterCreator(azure.Data),
	}
	if err := r.ReconcileClusters(ctx, clusterCreator, namespace, client); err != nil {
		return err
	}

	managedClusterCreator := []r.NamedAzureManagedClusterCreatorGetter{
		azureManagedClusterCreator(azure.Data),
	}
	if err := r.ReconcileAzureManagedClusters(ctx, managedClusterCreator, namespace, client); err != nil {
		return err
	}

	managedControlPlaneCreator := []r.NamedAzureManagedControlPlaneCreatorGetter{
		azureManageControlPlaneCreator(azure.Data),
	}
	if err := r.ReconcileAzureManagedControlPlanes(ctx, managedControlPlaneCreator, namespace, client); err != nil {
		return err
	}

	machinePoolCreator := []r.NamedMachinePoolCreatorGetter{}
	managedMachinePoolCreator := []r.NamedAzureManagedMachinePoolCreatorGetter{}
	for _, machinePool := range azure.Data.Bootstrap.Spec.CloudSpec.Azure.MachinePools {
		machinePoolCreator = append(machinePoolCreator, azureMachinePoolCreator(machinePool, azure.Data))
		managedMachinePoolCreator = append(managedMachinePoolCreator, azureManagedMachinePoolCreator(machinePool, azure.Data))
	}
	if err := r.ReconcileMachinePools(ctx, machinePoolCreator, namespace, client); err != nil {
		return err
	}
	if err := r.ReconcileAzureManagedMachinePools(ctx, managedMachinePoolCreator, namespace, client); err != nil {
		return err
	}

	return nil
}

func azureClusterIdentityCreator(data *resources.TemplateData) r.NamedAzureClusterIdentityCreatorGetter {
	return func() (string, r.AzureClusterIdentityCreator) {
		return data.Bootstrap.Spec.CloudSpec.Azure.ClusterIdentity.Name, func(c *azure.AzureClusterIdentity) (*azure.AzureClusterIdentity, error) {
			identity := data.Bootstrap.Spec.CloudSpec.Azure.ClusterIdentity

			c.Name = identity.Name
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = azure.AzureClusterIdentitySpec{
				Type:         azure.IdentityType(identity.Type),
				ResourceID:   identity.ResourceID,
				ClientID:     identity.ClientID,
				ClientSecret: identity.ClientSecret,
				TenantID:     identity.TenantID,
				AllowedNamespaces: &azure.AllowedNamespaces{
					NamespaceList: identity.AllowedNamespaces.NamespaceList,
					Selector:      identity.AllowedNamespaces.Selector,
				},
			}

			return c, nil
		}
	}
}

func azureClusterCreator(data *resources.TemplateData) r.NamedClusterCreatorGetter {
	return func() (string, r.ClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *clusterv1.Cluster) (*clusterv1.Cluster, error) {
			name := data.Bootstrap.Spec.ClusterName
			c.Name = name
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = clusterv1.ClusterSpec{
				ClusterNetwork: &clusterv1.ClusterNetwork{
					APIServerPort: data.Bootstrap.Spec.ClusterNetwork.APIServerPort,
					ServiceDomain: data.Bootstrap.Spec.ClusterNetwork.ServiceDomain,
				},
			}
			if data.Bootstrap.Spec.ClusterNetwork.Pods != nil {
				c.Spec.ClusterNetwork.Pods = &clusterv1.NetworkRanges{
					CIDRBlocks: data.Bootstrap.Spec.ClusterNetwork.Pods.CIDRBlocks,
				}
			}
			if data.Bootstrap.Spec.ClusterNetwork.Services != nil {
				c.Spec.ClusterNetwork.Services = &clusterv1.NetworkRanges{
					CIDRBlocks: data.Bootstrap.Spec.ClusterNetwork.Services.CIDRBlocks,
				}
			}
			c.Spec.ControlPlaneRef = &corev1.ObjectReference{
				Kind:       "AzureManagedControlPlane",
				Name:       name,
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			}
			c.Spec.InfrastructureRef = &corev1.ObjectReference{
				Kind:       "AzureManagedCluster",
				Name:       name,
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			}

			return c, nil
		}
	}
}

func azureManagedClusterCreator(data *resources.TemplateData) r.NamedAzureManagedClusterCreatorGetter {
	return func() (string, r.AzureManagedClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *azure.AzureManagedCluster) (*azure.AzureManagedCluster, error) {
			cluster := data.Bootstrap.Spec.CloudSpec.Azure.ManagedCluster

			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = azure.AzureManagedClusterSpec{
				ControlPlaneEndpoint: clusterv1.APIEndpoint(cluster.ControlPlaneEndpoint),
			}

			return c, nil
		}
	}
}

func azureManageControlPlaneCreator(data *resources.TemplateData) r.NamedAzureManagedControlPlaneCreatorGetter {
	return func() (string, r.AzureManagedControlPlaneCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *azure.AzureManagedControlPlane) (*azure.AzureManagedControlPlane, error) {
			controlPlane := data.Bootstrap.Spec.CloudSpec.Azure.ControlPlane

			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = azure.AzureManagedControlPlaneSpec{
				Version:           controlPlane.Version,
				ResourceGroupName: controlPlane.ResourceGroupName,
				SubscriptionID:    controlPlane.SubscriptionID,
				Location:          controlPlane.Location,
				SSHPublicKey:      controlPlane.SSHPublicKey,
				IdentityRef:       controlPlane.IdentityRef,
			}

			if len(data.Bootstrap.Spec.KubernetesVersion) > 0 {
				c.Spec.Version = data.Bootstrap.Spec.KubernetesVersion
			}

			return c, nil
		}
	}
}

func azureMachinePoolCreator(machinePool *bv1alpha1.AzureMachinePool, data *resources.TemplateData) r.NamedMachinePoolCreatorGetter {
	return func() (string, r.MachinePoolCreator) {
		return machinePool.Name, func(c *clusterapiexp.MachinePool) (*clusterapiexp.MachinePool, error) {
			c.Name = machinePool.Name
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = clusterapiexp.MachinePoolSpec{
				ClusterName: data.Bootstrap.Spec.ClusterName,
				Replicas:    machinePool.Replicas,
				Template: clusterv1.MachineTemplateSpec{
					Spec: clusterv1.MachineSpec{
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: resources.StrPtr(""),
						},
						ClusterName: data.Bootstrap.Spec.ClusterName,
						InfrastructureRef: corev1.ObjectReference{
							Kind:       "AzureManagedMachinePool",
							Name:       machinePool.Name,
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
						},
					},
				},
			}
			return c, nil
		}
	}
}

func azureManagedMachinePoolCreator(mp *bv1alpha1.AzureMachinePool, data *resources.TemplateData) r.NamedAzureManagedMachinePoolCreatorGetter {
	return func() (string, r.AzureManagedMachinePoolCreator) {
		return mp.Name, func(c *azure.AzureManagedMachinePool) (*azure.AzureManagedMachinePool, error) {
			c.Name = mp.Name
			c.Namespace = data.Bootstrap.Namespace
			c.Spec = azure.AzureManagedMachinePoolSpec{
				Name: &mp.Name,
				Mode: mp.Mode,
				SKU:  mp.SKU,
				ScaleSetPriority: mp.ScaleSetPriority,
				AvailabilityZones: mp.AvailabilityZones,
				OsDiskType: mp.OsDiskType,
				OSDiskSizeGB: mp.OSDiskSizeGB,
				MaxPods: mp.MaxPods,
				NodeLabels: mp.NodeLabels,
				AdditionalTags: azure.Tags(mp.AdditionalTags),
			}

			if mp.Scaling != nil {
				c.Spec.Scaling = &azure.ManagedMachinePoolScaling{
					MinSize: mp.Scaling.MinSize,
					MaxSize: mp.Scaling.MaxSize,
				}
			}

			for _, taint := range mp.Taints {
				c.Spec.Taints = append(c.Spec.Taints, azure.Taint{
					Effect: azure.TaintEffect(taint.Effect),
					Key:    taint.Key,
					Value:  taint.Value,
				})
			}

			return c, nil
		}
	}
}
