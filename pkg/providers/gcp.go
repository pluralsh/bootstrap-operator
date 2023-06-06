package providers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	cloudcontainer "cloud.google.com/go/container/apiv1"
	container "cloud.google.com/go/container/apiv1/containerpb"
	"github.com/pluralsh/polly/algorithms"
	"google.golang.org/api/option"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/version"
	"sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	gcpclusterapi "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	expgcpclusterapi "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
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
	gcpClient      *cloudcontainer.ClusterManagerClient
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

func (gcp *GCPProvider) MigrateCluster() (*ctrl.Result, error) {
	return nil, nil
}

func (gcp *GCPProvider) CheckCluster() (*ctrl.Result, error) {
	cluster := new(clusterapi.Cluster)
	if err := gcp.Data.Client.Get(gcp.Data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: gcp.Data.Namespace, Name: gcp.Data.Bootstrap.Spec.ClusterName}, cluster); err != nil {
		return nil, err
	}

	if err := gcp.updateClusterStatus(cluster.Status); err != nil {
		return nil, err
	}

	status, err := gcp.getClusterReadyStatus(cluster.Status)
	if err != nil {
		return nil, err
	}

	if *status != corev1.ConditionTrue {
		gcp.Data.Log.WithName("GCP provider").Info("Waiting for cluster to become ready...")
		return &ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	gcp.Data.Log.WithName("GCP provider").Info("Cluster provisioned. Running preflight checks...")
	if res, err := gcp.runClusterPreflightChecks(); err != nil || res != nil {
		return res, err
	}

	gcp.Data.Log.WithName("GCP provider").Info("Updating cluster object Ready status to true")
	return nil, gcp.setStatusReady()
}

func (gcp *GCPProvider) getClusterReadyStatus(status clusterapi.ClusterStatus) (*corev1.ConditionStatus, error) {
	matchingConditions := algorithms.Filter(status.Conditions, func(condition clusterv1.Condition) bool {
		return condition.Type == clusterv1.ReadyCondition
	})

	if len(matchingConditions) != 1 {
		return nil, fmt.Errorf("could not find cluster Ready condition")
	}

	return &matchingConditions[0].Status, nil
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
	cloudSpec := gcp.Data.Bootstrap.Spec.CloudSpec.GCP

	clusterCreator := []reconciling.NamedClusterCreatorGetter{clusterCreator(gcp.Data)}
	managedClusterCreator := []reconciling.NamedGCPManagedClusterCreatorGetter{gcpManagedClusterCreator(gcp.Data)}
	managedControlPlaneCreator := []reconciling.NamedGCPManagedControlPlaneCreatorGetter{gcpManagedControlPlaneCreator(gcp.Data)}
	machinePoolCreator := []reconciling.NamedMachinePoolCreatorGetter{gcpMachinePoolCreator(gcp.Data)}
	managedMachinePoolCreator := []reconciling.NamedGCPManagedMachinePoolCreatorGetter{gcpManagedMachinePoolCreator(gcp.Data)}

	if err := reconciling.ReconcileClusters(gcp.Data.Ctx, clusterCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	if err := reconciling.ReconcileGCPManagedClusters(gcp.Data.Ctx, managedClusterCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	if err := reconciling.ReconcileGCPManagedControlPlanes(gcp.Data.Ctx, managedControlPlaneCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	// Do not create any machine pools if GKE autopilot is enabled.
	if cloudSpec.ControlPlane != nil && cloudSpec.ControlPlane.EnableAutopilot {
		return nil
	}

	if err := reconciling.ReconcileMachinePools(gcp.Data.Ctx, machinePoolCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	if err := reconciling.ReconcileGCPManagedMachinePools(gcp.Data.Ctx, managedMachinePoolCreator, gcp.Data.Namespace, gcp.Data.Client); err != nil {
		return err
	}

	return nil
}

func GetGCPProvider(data *resources.TemplateData) (*GCPProvider, error) {
	ctx := context.Background()
	cloudSpec := data.Bootstrap.Spec.CloudSpec.GCP

	var secret corev1.Secret
	if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: data.Namespace, Name: cloudSpec.CredentialsRef.Name}, &secret); err != nil {
		return nil, err
	}

	credentials := strings.TrimSpace(string(secret.Data[cloudSpec.CredentialsRef.Key]))
	gcpClient, err := cloudcontainer.NewClusterManagerClient(ctx, defaultClientOptions(credentials)...)
	if err != nil {
		return nil, fmt.Errorf("could not create gcp client: %s", err)
	}

	data.Log.WithName("GCP provider").Info("Create GCP provider")
	return &GCPProvider{
		Data:           data,
		Credentials:    credentials,
		Region:         cloudSpec.Region,
		version:        cloudSpec.Version,
		fetchConfigUrl: cloudSpec.FetchConfigURL,
		gcpClient:      gcpClient,
	}, nil
}

func clusterCreator(data *resources.TemplateData) reconciling.NamedClusterCreatorGetter {
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
		return data.Bootstrap.Spec.ClusterName, func(c *expgcpclusterapi.GCPManagedCluster) (*expgcpclusterapi.GCPManagedCluster, error) {
			subnets := make([]gcpclusterapi.SubnetSpec, 0)
			for _, subnet := range data.Bootstrap.Spec.CloudSpec.GCP.Cluster.Network.Subnets {
				subnets = append(subnets, gcpclusterapi.SubnetSpec{
					Name:                subnet.Name,
					CidrBlock:           subnet.CidrBlock,
					SecondaryCidrBlocks: subnet.SecondaryCidrBlocks,
					Region:              data.Bootstrap.Spec.CloudSpec.GCP.Region,
				})
			}

			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Namespace
			c.Spec = expgcpclusterapi.GCPManagedClusterSpec{
				Project: data.Bootstrap.Spec.CloudSpec.GCP.Project,
				Region:  data.Bootstrap.Spec.CloudSpec.GCP.Region,
				Network: v1beta1.NetworkSpec{
					Name:                  data.Bootstrap.Spec.CloudSpec.GCP.Cluster.Network.Name,
					AutoCreateSubnetworks: data.Bootstrap.Spec.CloudSpec.GCP.Cluster.Network.AutoCreateSubnetworks,
					Subnets:               subnets,
				},
			}

			return c, nil
		}
	}
}

func gcpManagedControlPlaneCreator(data *resources.TemplateData) reconciling.NamedGCPManagedControlPlaneCreatorGetter {
	return func() (string, reconciling.GCPManagedControlPlaneCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane"), func(c *expgcpclusterapi.GCPManagedControlPlane) (*expgcpclusterapi.GCPManagedControlPlane, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane")
			c.Namespace = data.Namespace
			c.Spec = expgcpclusterapi.GCPManagedControlPlaneSpec{
				ClusterName: data.Bootstrap.Spec.ClusterName,
				Project:     data.Bootstrap.Spec.CloudSpec.GCP.Project,
				Location:    data.Bootstrap.Spec.CloudSpec.GCP.Region,
			}

			if len(data.Bootstrap.Spec.KubernetesVersion) > 0 {
				c.Spec.ControlPlaneVersion = resources.StrPtr(data.Bootstrap.Spec.KubernetesVersion)
			}

			if data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane == nil {
				return c, nil
			}

			c.Spec.EnableAutopilot = data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane.EnableAutopilot
			c.Spec.EnableWorkloadIdentity = data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane.EnableWorkloadIdentity
			if data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane.ReleaseChannel != nil {
				var channel expgcpclusterapi.ReleaseChannel

				switch *data.Bootstrap.Spec.CloudSpec.GCP.ControlPlane.ReleaseChannel {
				case bv1alpha1.Rapid:
					channel = expgcpclusterapi.Rapid
				case bv1alpha1.Regular:
					channel = expgcpclusterapi.Regular
				case bv1alpha1.Stable:
					channel = expgcpclusterapi.Stable
				}

				c.Spec.ReleaseChannel = &channel
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

func gcpManagedMachinePoolCreator(data *resources.TemplateData) reconciling.NamedGCPManagedMachinePoolCreatorGetter {
	return func() (string, reconciling.GCPManagedMachinePoolCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0"), func(c *expgcpclusterapi.GCPManagedMachinePool) (*expgcpclusterapi.GCPManagedMachinePool, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "pool-0")
			c.Namespace = data.Namespace
			c.Spec = expgcpclusterapi.GCPManagedMachinePoolSpec{}

			if data.Bootstrap.Spec.CloudSpec.GCP.MachinePool.Scaling != nil {
				c.Spec.Scaling = &expgcpclusterapi.NodePoolAutoScaling{
					MinCount: data.Bootstrap.Spec.CloudSpec.GCP.MachinePool.Scaling.MinCount,
					MaxCount: data.Bootstrap.Spec.CloudSpec.GCP.MachinePool.Scaling.MinCount,
				}
			}

			return c, nil
		}
	}
}

func (gcp *GCPProvider) runClusterPreflightChecks() (*ctrl.Result, error) {
	running, err := gcp.isClusterRunning()
	if err != nil {
		return nil, err
	}

	if !running {
		gcp.Data.Log.WithName("GCP provider").Info("Waiting for cluster to be in running state...")
		return &ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}

	return nil, nil
}

func (gcp *GCPProvider) isClusterRunning() (bool, error) {
	ctx := context.Background()
	cluster, err := gcp.gcpClient.GetCluster(ctx, &container.GetClusterRequest{Name: gcp.clusterName()})
	if err != nil {
		return false, err
	}

	return cluster.Status == container.Cluster_RUNNING, nil
}

func (gcp *GCPProvider) clusterName() string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		gcp.Data.Bootstrap.Spec.CloudSpec.GCP.Project,
		gcp.Data.Bootstrap.Spec.CloudSpec.GCP.Region,
		gcp.Data.Bootstrap.Spec.ClusterName,
	)
}

func defaultClientOptions(credentials string) []option.ClientOption {
	return []option.ClientOption{
		option.WithUserAgent(fmt.Sprintf("gcp.cluster.x-k8s.io/%s", version.Get())),
		option.WithCredentialsJSON([]byte(credentials)),
	}
}
