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
	Location       string
	ClientSecret   string
	ClientID       string
	TenantID       string
	SubscriptionID string
	GitHubToken    string
	version        string
	fetchConfigUrl string
}

const (
	azureSecretName                = "azure-credentials"
	azureClusterIdentityName       = "azure-cluster-identity"
	azureClusterIdentitySecretName = "azure-cluster-identity-secret"
)

func GetAzureProvider(data *resources.TemplateData) (*AzureProvider, error) {
	spec := data.Bootstrap.Spec.CloudSpec.Azure

	// TODO: Right now it assumes that same secret is used for all refs, that's why client is retrieved only once.
	// It should be changed.
	var secret corev1.Secret
	if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: data.Namespace, Name: spec.ClientSecretRef.Name}, &secret); err != nil {
		return nil, err
	}

	clientSecret := strings.TrimSpace(string(secret.Data[spec.ClientSecretRef.Key]))
	clientID := strings.TrimSpace(string(secret.Data[spec.ClientIDRef.Key]))
	tenantID := strings.TrimSpace(string(secret.Data[spec.TenantIDRef.Key]))
	subscriptionID := strings.TrimSpace(string(secret.Data[spec.SubscriptionIDRef.Key]))
	gitHubToken := strings.TrimSpace(string(secret.Data[data.Bootstrap.Spec.GitHubSecretRef.Key]))
	data.Log.Named("Azure provider").Info("Create Azure provider")
	return &AzureProvider{
		Data:           data,
		version:        spec.Version,
		fetchConfigUrl: spec.FetchConfigURL,
		Location:       spec.ControlPlane.Location,
		ClientSecret:   clientSecret,
		ClientID:       clientID,
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
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
			"AZURE_LOCATION":                          []byte(azure.Location),
			"AZURE_CLIENT_SECRET":                     []byte(azure.ClientSecret),
			"AZURE_CLIENT_ID":                         []byte(azure.ClientID),
			"AZURE_TENANT_ID":                         []byte(azure.TenantID),
			"AZURE_SUBSCRIPTION_ID":                   []byte(azure.SubscriptionID),
			"AZURE_CLUSTER_IDENTITY_SECRET_NAME":      []byte(azureClusterIdentitySecretName),
			"AZURE_CLUSTER_IDENTITY_SECRET_NAMESPACE": []byte(azure.Data.Namespace),
			"EXP_MACHINE_POOL":                        []byte("true"),
			"GITHUB_TOKEN":                            []byte(azure.GitHubToken),
		},
	}
	if err := azure.Data.Client.Create(azure.Data.Ctx, &secret); err != nil {
		return err
	}
	return nil
}

func (azure *AzureProvider) createClusterIdentitySecret() error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: azure.Data.Namespace,
			Name:      azureClusterIdentitySecretName,
		},
		Data: map[string][]byte{
			"clientSecret": []byte(azure.ClientSecret),
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
		if err := azure.createClusterIdentitySecret(); err != nil {
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
		azureClusterIdentityCreator(azure),
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

	machinePoolCreator := []reconciling.NamedMachinePoolCreatorGetter{
		azureMachinePoolCreator(azure.Data),
	}
	if err := reconciling.ReconcileMachinePools(azure.Data.Ctx, machinePoolCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	managedMachinePoolCreator := []reconciling.NamedAzureManagedMachinePoolCreatorGetter{
		azureManagedMachinePoolCreator(azure.Data),
		azureManagedUserMachinePoolCreator(azure.Data),
	}
	if err := reconciling.ReconcileAzureManagedMachinePools(azure.Data.Ctx, managedMachinePoolCreator, azure.Data.Namespace, azure.Data.Client); err != nil {
		return err
	}

	return nil
}

func azureClusterIdentityCreator(azureProvider *AzureProvider) reconciling.NamedAzureClusterIdentityCreatorGetter {
	return func() (string, reconciling.AzureClusterIdentityCreator) {
		return azureClusterIdentityName, func(c *azure.AzureClusterIdentity) (*azure.AzureClusterIdentity, error) {
			c.Name = azureClusterIdentityName
			c.Namespace = azureProvider.Data.Namespace
			c.Spec = azure.AzureClusterIdentitySpec{
				Type:     azure.ServicePrincipal,
				ClientID: azureProvider.ClientID,
				ClientSecret: corev1.SecretReference{
					Name:      azureSecretName,
					Namespace: azureProvider.Data.Namespace,
				},
				TenantID: azureProvider.TenantID,
			}

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
				Replicas:    resources.Int32(2), // TODO: Change it.
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
			// TODO: Check spec
			c.Spec = azure.AzureManagedMachinePoolSpec{
				Mode: "User",
				SKU:  "Standard_D2s_v3",
			}

			return c, nil
		}
	}
}

func azureManagedUserMachinePoolCreator(data *resources.TemplateData) reconciling.NamedAzureManagedMachinePoolCreatorGetter {
	return func() (string, reconciling.AzureManagedMachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *azure.AzureManagedMachinePool) (*azure.AzureManagedMachinePool, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Namespace = data.Namespace
			// TODO: Check spec change to system
			c.Spec = azure.AzureManagedMachinePoolSpec{
				Mode: "User",
				SKU:  "Standard_D2s_v3",
			}

			return c, nil
		}
	}
}
