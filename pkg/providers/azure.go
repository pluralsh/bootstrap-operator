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
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
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

	var secret corev1.Secret
	if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{
		Namespace: data.Namespace, Name: spec.GitHubSecretRef.Name}, &secret); err != nil {
		return nil, err
	}
	gitHubToken := strings.TrimSpace(string(secret.Data[spec.GitHubSecretRef.Key]))

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
			Namespace: azure.Data.Namespace,
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
		Namespace: azure.Data.Namespace,
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
	clusterIdentityCreator := []reconciling.NamedAzureClusterIdentityCreatorGetter{
		azureClusterIdentityCreator(azure.Data),
	}
	if err := reconciling.ReconcileAzureClusterIdentitys(azure.Data.Ctx, clusterIdentityCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	clusterCreator := []reconciling.NamedClusterCreatorGetter{
		azureClusterCreator(azure.Data),
	}
	if err := reconciling.ReconcileClusters(azure.Data.Ctx, clusterCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	managedClusterCreator := []reconciling.NamedAzureManagedClusterCreatorGetter{
		azureManagedClusterCreator(azure.Data),
	}
	if err := reconciling.ReconcileAzureManagedClusters(azure.Data.Ctx, managedClusterCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	managedControlPlaneCreator := []reconciling.NamedAzureManagedControlPlaneCreatorGetter{
		azureManageControlPlaneCreator(azure.Data),
	}
	if err := reconciling.ReconcileAzureManagedControlPlanes(azure.Data.Ctx, managedControlPlaneCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	// TODO: At the moment only one machine pool with one respective managed machine pool will be created.
	// In the future it should be possible to specify multiple machine pools at the cloud spec level.
	machinePoolCreator := []reconciling.NamedMachinePoolCreatorGetter{
		azureMachinePoolCreator(azure.Data),
	}
	if err := reconciling.ReconcileMachinePools(azure.Data.Ctx, machinePoolCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	managedMachinePoolCreator := []reconciling.NamedAzureManagedMachinePoolCreatorGetter{
		azureManagedMachinePoolCreator(azure.Data),
	}
	if err := reconciling.ReconcileAzureManagedMachinePools(azure.Data.Ctx, managedMachinePoolCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	return nil
}

func azureClusterIdentityCreator(data *resources.TemplateData) reconciling.NamedAzureClusterIdentityCreatorGetter {
	return func() (string, reconciling.AzureClusterIdentityCreator) {
		return data.Bootstrap.Spec.CloudSpec.Azure.ClusterIdentity.Name, func(c *azure.AzureClusterIdentity) (*azure.AzureClusterIdentity, error) {
			c.Name = data.Bootstrap.Spec.CloudSpec.Azure.ClusterIdentity.Name
			c.Namespace = data.Namespace
			c.Spec = data.Bootstrap.Spec.CloudSpec.Azure.ClusterIdentity.AzureClusterIdentitySpec

			return c, nil
		}
	}
}

func azureClusterCreator(data *resources.TemplateData) reconciling.NamedClusterCreatorGetter {
	return func() (string, reconciling.ClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *clusterv1.Cluster) (*clusterv1.Cluster, error) {
			name := data.Bootstrap.Spec.ClusterName
			c.Name = name
			c.Namespace = data.Namespace
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
				Name:       fmt.Sprintf("%s-%s", name, "control-plane"),
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

func azureManagedClusterCreator(data *resources.TemplateData) reconciling.NamedAzureManagedClusterCreatorGetter {
	return func() (string, reconciling.AzureManagedClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *azure.AzureManagedCluster) (*azure.AzureManagedCluster, error) {
			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Namespace
			c.Spec = *data.Bootstrap.Spec.CloudSpec.Azure.ManagedCluster

			return c, nil
		}
	}
}

func azureManageControlPlaneCreator(data *resources.TemplateData) reconciling.NamedAzureManagedControlPlaneCreatorGetter {
	return func() (string, reconciling.AzureManagedControlPlaneCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane"), func(c *azure.AzureManagedControlPlane) (*azure.AzureManagedControlPlane, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane")
			c.Namespace = data.Namespace
			c.Spec = *data.Bootstrap.Spec.CloudSpec.Azure.ControlPlane

			if len(data.Bootstrap.Spec.KubernetesVersion) > 0 {
				c.Spec.Version = data.Bootstrap.Spec.KubernetesVersion
			}

			return c, nil
		}
	}
}

func azureMachinePoolCreator(data *resources.TemplateData) reconciling.NamedMachinePoolCreatorGetter {
	return func() (string, reconciling.MachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *clusterapiexp.MachinePool) (*clusterapiexp.MachinePool, error) {
			name := fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Name = name
			c.Namespace = data.Namespace
			c.Spec = clusterapiexp.MachinePoolSpec{
				ClusterName: data.Bootstrap.Spec.ClusterName,
				Replicas:    resources.Int32(3),
				Template: clusterv1.MachineTemplateSpec{
					Spec: clusterv1.MachineSpec{
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: resources.StrPtr(""),
						},
						ClusterName: data.Bootstrap.Spec.ClusterName,
						InfrastructureRef: corev1.ObjectReference{
							Kind:       "AzureManagedMachinePool",
							Name:       name,
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
						},
					},
				},
			}
			return c, nil
		}
	}
}

func azureManagedMachinePoolCreator(data *resources.TemplateData) reconciling.NamedAzureManagedMachinePoolCreatorGetter {
	return func() (string, reconciling.AzureManagedMachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *azure.AzureManagedMachinePool) (*azure.AzureManagedMachinePool, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Namespace = data.Namespace
			c.Spec = azure.AzureManagedMachinePoolSpec{
				Mode: "System",
				SKU:  "Standard_D2s_v3",
			}

			return c, nil
		}
	}
}
