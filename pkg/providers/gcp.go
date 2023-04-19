package providers

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gcpclusterapi "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterapi "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

type GCPProvider struct {
	Data           *resources.TemplateData
	Credentials    string
	Region         string
	version        string
	fetchConfigUrl string
}

const (
	gcpSecretName = "gcp-credentials"
)

func (gcp *GCPProvider) Name() string {
	return "gcp"
}

func (gcp *GCPProvider) Version() string {
	return gcp.version
}

func (gcp *GCPProvider) FetchConfigURL() string {
	return gcp.fetchConfigUrl
}

func (gcp *GCPProvider) createCredentialSecret() error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gcp.Data.Namespace,
			Name:      gcpSecretName,
		},
		Data: map[string][]byte{
			"GCP_REGION":                 []byte(gcp.Region),
			"EXP_CAPG_GKE":               []byte("true"),
			"EXP_MACHINE_POOL":           []byte("true"),
			"GCP_B64ENCODED_CREDENTIALS": []byte(base64.StdEncoding.EncodeToString([]byte(gcp.Credentials))),
		},
	}
	if err := gcp.Data.Client.Create(gcp.Data.Ctx, &secret); err != nil {
		return err
	}
	return nil
}

func (gcp *GCPProvider) Init() (*ctrl.Result, error) {
	if gcp.Data.Bootstrap.Status.ProviderStatus.Phase != bv1alpha1.Creating {
		if err := gcp.updateProviderStatus(bv1alpha1.Creating, "init GCP provider", false); err != nil {
			return nil, err
		}
		if err := gcp.createCredentialSecret(); err != nil {
			if err := gcp.updateProviderStatus(bv1alpha1.Error, err.Error(), false); err != nil {
				return nil, err
			}
		}
		if err := gcp.updateProviderStatus(bv1alpha1.Running, "GCP provider ready", true); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (gcp *GCPProvider) updateProviderStatus(phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(gcp.Data.Ctx, gcp.Data.Client, gcp.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
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

func (gcp *GCPProvider) Secret() string {
	return gcpSecretName
}

func (gcp *GCPProvider) CheckCluster() (*ctrl.Result, error) {
	var cluster clusterapi.Cluster
	if err := gcp.Data.Client.Get(gcp.Data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: gcp.Data.Namespace, Name: gcp.Data.Bootstrap.Spec.ClusterName}, &cluster); err != nil {
		return nil, err
	}
	if err := gcp.updateClusterStatus(cluster.Status); err != nil {
		return nil, err
	}
	for _, cond := range cluster.Status.Conditions {
		if cond.Type == clusterv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			return nil, gcp.setStatusReady()
		}
	}
	gcp.Data.Log.Named("GCP provider").Info("Waiting for GCP cluster to become ready")
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (gcp *GCPProvider) setStatusReady() error {
	return helper.UpdateBootstrapStatus(gcp.Data.Ctx, gcp.Data.Client, gcp.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiClusterStatus == nil {
			c.Status.CapiClusterStatus = &bv1alpha1.ClusterStatus{}
		}
		c.Status.CapiClusterStatus.Ready = true
	})
}

func (gcp *GCPProvider) updateClusterStatus(status clusterapi.ClusterStatus) error {
	err := helper.UpdateBootstrapStatus(gcp.Data.Ctx, gcp.Data.Client, gcp.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
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

func (gcp *GCPProvider) ReconcileCluster() error {
	clusterCreator := []reconciling.NamedClusterCreatorGetter{
		gcpClusterCreator(gcp.Data),
	}

	if err := reconciling.ReconcileClusters(gcp.Data.Ctx, clusterCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}
	managecClusterCreator := []reconciling.NamedGCPManagedClusterCreatorGetter{
		gcpManagedClusterCreator(gcp.Data),
	}
	if err := reconciling.ReconcileGCPManagedClusters(gcp.Data.Ctx, managecClusterCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}
	managecControlPlaneCreator := []reconciling.NamedGCPManagedControlPlaneCreatorGetter{
		gcpManageControlPlaneCreator(gcp.Data),
	}
	if err := reconciling.ReconcileGCPManagedControlPlanes(gcp.Data.Ctx, managecControlPlaneCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}
	machinePoolCreator := []reconciling.NamedMachinePoolCreatorGetter{
		gcpMachinePoolCreator(gcp.Data),
	}
	if err := reconciling.ReconcileMachinePools(gcp.Data.Ctx, machinePoolCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}
	managecMachinePoolCreator := []reconciling.NamedGCPManagedMachinePoolCreatorGetter{
		gcpGCPManagedMachinePoolCreator(gcp.Data),
	}
	if err := reconciling.ReconcileGCPManagedMachinePools(gcp.Data.Ctx, managecMachinePoolCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	return nil
}

func GetGCPProvider(data *resources.TemplateData) (*GCPProvider, error) {
	spec := data.Bootstrap.Spec.CloudSpec.GCP

	var secret corev1.Secret
	if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: data.Namespace, Name: spec.CredentialsRef.Name}, &secret); err != nil {
		return nil, err
	}

	credentials := strings.TrimSpace(string(secret.Data[spec.CredentialsRef.Key]))
	data.Log.Named("GCP provider").Info("Create GCP provider")
	return &GCPProvider{
		Data:           data,
		Credentials:    credentials,
		Region:         spec.ManagedCluster.Region,
		version:        spec.Version,
		fetchConfigUrl: spec.FetchConfigURL,
	}, nil
}

func gcpClusterCreator(data *resources.TemplateData) reconciling.NamedClusterCreatorGetter {
	return func() (string, reconciling.ClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *clusterapi.Cluster) (*clusterapi.Cluster, error) {
			name := data.Bootstrap.Spec.ClusterName
			c.Name = name
			c.Namespace = data.Namespace
			c.Spec = clusterapi.ClusterSpec{
				ClusterNetwork: &clusterapi.ClusterNetwork{
					APIServerPort: data.Bootstrap.Spec.ClusterNetwork.APIServerPort,
					ServiceDomain: data.Bootstrap.Spec.ClusterNetwork.ServiceDomain,
				},
			}
			if data.Bootstrap.Spec.ClusterNetwork.Pods != nil {
				c.Spec.ClusterNetwork.Pods = &clusterapi.NetworkRanges{
					CIDRBlocks: data.Bootstrap.Spec.ClusterNetwork.Pods.CIDRBlocks,
				}
			}
			if data.Bootstrap.Spec.ClusterNetwork.Services != nil {
				c.Spec.ClusterNetwork.Services = &clusterapi.NetworkRanges{
					CIDRBlocks: data.Bootstrap.Spec.ClusterNetwork.Services.CIDRBlocks,
				}
			}
			c.Spec.ControlPlaneRef = &corev1.ObjectReference{
				Kind:       "GCPManagedControlPlane",
				Name:       fmt.Sprintf("%s-%s", name, "control-plane"),
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			}
			c.Spec.InfrastructureRef = &corev1.ObjectReference{
				Kind:       "GCPManagedCluster",
				Name:       name,
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			}

			return c, nil
		}
	}
}

func gcpManagedClusterCreator(data *resources.TemplateData) reconciling.NamedGCPManagedClusterCreatorGetter {
	return func() (string, reconciling.GCPManagedClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *gcpclusterapi.GCPManagedCluster) (*gcpclusterapi.GCPManagedCluster, error) {
			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Namespace
			c.Spec = *data.Bootstrap.Spec.CloudSpec.GCP.ManagedCluster

			return c, nil
		}
	}
}

func gcpManageControlPlaneCreator(data *resources.TemplateData) reconciling.NamedGCPManagedControlPlaneCreatorGetter {
	return func() (string, reconciling.GCPManagedControlPlaneCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane"), func(c *gcpclusterapi.GCPManagedControlPlane) (*gcpclusterapi.GCPManagedControlPlane, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane")
			c.Namespace = data.Namespace
			c.Spec = *data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane
			c.Spec.ClusterName = data.Bootstrap.Spec.ClusterName

			if len(data.Bootstrap.Spec.KubernetesVersion) > 0 {
				c.Spec.ControlPlaneVersion = resources.StrPtr(data.Bootstrap.Spec.KubernetesVersion)
			}

			return c, nil
		}
	}
}

func gcpMachinePoolCreator(data *resources.TemplateData) reconciling.NamedMachinePoolCreatorGetter {
	return func() (string, reconciling.MachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *clusterapiexp.MachinePool) (*clusterapiexp.MachinePool, error) {
			name := fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Name = name
			c.Namespace = data.Namespace
			c.Spec = clusterapiexp.MachinePoolSpec{
				ClusterName: data.Bootstrap.Spec.ClusterName,
				Replicas:    resources.Int32(data.Bootstrap.Spec.CloudSpec.GCP.MachinePool.Replicas),
				Template: clusterapi.MachineTemplateSpec{
					Spec: clusterapi.MachineSpec{
						Bootstrap: clusterapi.Bootstrap{
							DataSecretName: resources.StrPtr(""),
						},
						ClusterName: data.Bootstrap.Spec.ClusterName,
						InfrastructureRef: corev1.ObjectReference{
							Kind:       "GCPManagedMachinePool",
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

func gcpGCPManagedMachinePoolCreator(data *resources.TemplateData) reconciling.NamedGCPManagedMachinePoolCreatorGetter {
	return func() (string, reconciling.GCPManagedMachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *gcpclusterapi.GCPManagedMachinePool) (*gcpclusterapi.GCPManagedMachinePool, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Namespace = data.Namespace
			// TODO: Check spec
			c.Spec = gcpclusterapi.GCPManagedMachinePoolSpec{}

			return c, nil
		}
	}
}
