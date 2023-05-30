package providers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/yaml"

	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	awsinfrastructure "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/service"
	creds "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/credentials"
	awscontrolplane "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	awsmachinepool "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
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

type AWSProvider struct {
	Data            *resources.TemplateData
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Credentials     string
	Region          string
	version         string
	fetchConfigURL  string
}

const (
	awsSecretName = "aws-credentials"
)

func (aws *AWSProvider) Name() string {
	return "aws"
}

func (aws *AWSProvider) Version() string {
	return aws.version
}

func (aws *AWSProvider) FetchConfigURL() string {
	return aws.fetchConfigURL
}

func (aws *AWSProvider) createCredentialSecret() error {

	os.Setenv("AWS_ACCESS_KEY_ID", aws.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", aws.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", aws.SessionToken)

	t := bootstrap.NewTemplate()
	t.Spec.Region = aws.Region

	config := awsv1.
		NewConfig().
		WithRegion(aws.Region).
		WithCredentials(credentials.NewStaticCredentials(aws.AccessKeyID, aws.SecretAccessKey, aws.SessionToken)).
		WithMaxRetries(3)

	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}
	cfnSvc := cloudformation.NewService(cfn.New(sess))
	cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *t.RenderCloudFormation(), t.Spec.StackTags)
	awsCreds, err := creds.NewAWSCredentialFromDefaultChain(aws.Region)
	if err != nil {
		return err
	}
	credentials, err := awsCreds.RenderBase64EncodedAWSDefaultProfile()
	if err != nil {
		return err
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: aws.Data.Namespace,
			Name:      awsSecretName,
		},
		Data: map[string][]byte{
			"CAPA_EKS_ADD_ROLES":       []byte("true"),
			"CAPA_EKS_IAM":             []byte("true"),
			"AWS_REGION":               []byte(aws.Region),
			"EXP_MACHINE_POOL":         []byte("true"),
			"EXP_EXTERNAL_RESOURCE_GC": []byte("true"),
		},
	}

	if aws.Data.Bootstrap.Spec.BootstrapMode {
		secret.Data["AWS_B64ENCODED_CREDENTIALS"] = []byte(credentials)
	} else {
		secret.Data["AWS_B64ENCODED_CREDENTIALS"] = []byte("")
		secret.Data["AWS_CONTROLLER_IAM_ROLE"] = []byte(fmt.Sprintf("arn:aws:iam::%s:role/capa-controller-manager", aws.AccountID))
	}

	if err := aws.Data.Client.Create(aws.Data.Ctx, &secret); err != nil {
		return err
	}
	return nil
}

func (aws *AWSProvider) Init() (*ctrl.Result, error) {
	if aws.Data.Bootstrap.Status.ProviderStatus.Phase != bv1alpha1.Creating {
		if err := aws.updateProviderStatus(bv1alpha1.Creating, "init AWS provider", false); err != nil {
			return nil, err
		}
		if err := aws.createCredentialSecret(); err != nil {
			if err := aws.updateProviderStatus(bv1alpha1.Error, err.Error(), false); err != nil {
				return nil, err
			}
		}
		if err := aws.updateProviderStatus(bv1alpha1.Running, "AWS provider ready", true); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (aws *AWSProvider) updateProviderStatus(phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(aws.Data.Ctx, aws.Data.Client, aws.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
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

func (aws *AWSProvider) Secret() string {
	return awsSecretName
}

func (aws *AWSProvider) CheckCluster() (*ctrl.Result, error) {
	var cluster clusterapi.Cluster
	if err := aws.Data.Client.Get(aws.Data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: aws.Data.Namespace, Name: aws.Data.Bootstrap.Spec.ClusterName}, &cluster); err != nil {
		return nil, err
	}
	if err := aws.updateClusterStatus(cluster.Status); err != nil {
		return nil, err
	}
	for _, cond := range cluster.Status.Conditions {
		if cond.Type == clusterv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			if err := aws.postInstall(); err != nil {
				return nil, err
			}
			return nil, aws.setStatusReady()
		}
	}
	aws.Data.Log.Named("AWS provider").Info("Waiting for AWS cluster to become ready")
	return &ctrl.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (aws *AWSProvider) setStatusReady() error {
	return helper.UpdateBootstrapStatus(aws.Data.Ctx, aws.Data.Client, aws.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiClusterStatus == nil {
			c.Status.CapiClusterStatus = &bv1alpha1.ClusterStatus{}
		}
		c.Status.CapiClusterStatus.Ready = true
	})
}

func (aws *AWSProvider) updateClusterStatus(status clusterapi.ClusterStatus) error {
	err := helper.UpdateBootstrapStatus(aws.Data.Ctx, aws.Data.Client, aws.Data.Bootstrap, func(c *bv1alpha1.Bootstrap) {
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

func (aws *AWSProvider) ReconcileCluster() error {

	clusterCreator := []reconciling.NamedClusterCreatorGetter{
		awsClusterCreator(aws.Data),
	}

	if err := reconciling.ReconcileClusters(aws.Data.Ctx, clusterCreator, aws.Data.Namespace, aws.Data.Client); err != nil {
		return err
	}
	managecClusterCreator := []reconciling.NamedAWSManagedClusterCreatorGetter{
		awsManagedClusterCreator(aws.Data),
	}
	if err := reconciling.ReconcileAWSManagedClusters(aws.Data.Ctx, managecClusterCreator, aws.Data.Namespace, aws.Data.Client); err != nil {
		return err
	}
	managecControlPlaneCreator := []reconciling.NamedAWSManagedControlPlaneCreatorGetter{
		awsManageControlPlaneCreator(aws.Data),
	}
	if err := reconciling.ReconcileAWSManagedControlPlanes(aws.Data.Ctx, managecControlPlaneCreator, aws.Data.Namespace, aws.Data.Client); err != nil {
		return err
	}

	machinePoolCreator := []reconciling.NamedMachinePoolCreatorGetter{}
	managecMachinePoolCreator := []reconciling.NamedAWSManagedMachinePoolCreatorGetter{}

	for _, mp := range aws.Data.Bootstrap.Spec.CloudSpec.AWS.MachinePools {
		machinePoolCreator = append(machinePoolCreator, awsMachinePoolCreator(mp, aws.Data.Namespace, aws.Data.Bootstrap.Spec.ClusterName))
		managecMachinePoolCreator = append(managecMachinePoolCreator, awsAWSManagedMachinePoolCreator(mp, aws.Data.Namespace))
	}

	if err := reconciling.ReconcileMachinePools(aws.Data.Ctx, machinePoolCreator, aws.Data.Namespace, aws.Data.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileAWSManagedMachinePools(aws.Data.Ctx, managecMachinePoolCreator, aws.Data.Namespace, aws.Data.Client); err != nil {
		return err
	}

	return nil
}

func GetAWSProvider(data *resources.TemplateData) (*AWSProvider, error) {
	spec := data.Bootstrap.Spec.CloudSpec.AWS

	var secret corev1.Secret
	if err := data.Client.Get(data.Ctx, ctrlruntimeclient.ObjectKey{Namespace: data.Namespace, Name: spec.AccessKeyIDRef.Name}, &secret); err != nil {
		return nil, err
	}
	accessKeyID := strings.TrimSpace(string(secret.Data[spec.AccessKeyIDRef.Key]))
	secretAccessKey := strings.TrimSpace(string(secret.Data[spec.SecretAccessKeyRef.Key]))
	sessionToken := strings.TrimSpace(string(secret.Data[spec.SessionTokenRef.Key]))
	accountID := strings.TrimSpace(string(secret.Data[spec.AWSAccountIDRef.Key]))
	data.Log.Named("AWS provider").Info("Create AWS provider")
	return &AWSProvider{
		Data:            data,
		AccountID:       accountID,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
		Region:          spec.Region,
		version:         spec.Version,
		fetchConfigURL:  spec.FetchConfigURL,
	}, nil
}

func awsClusterCreator(data *resources.TemplateData) reconciling.NamedClusterCreatorGetter {
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
				Kind:       "AWSManagedControlPlane",
				Name:       fmt.Sprintf("%s-%s", name, "control-plane"),
				APIVersion: "controlplane.cluster.x-k8s.io/v1beta2",
			}
			c.Spec.InfrastructureRef = &corev1.ObjectReference{
				Kind:       "AWSManagedCluster",
				Name:       name,
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
			}

			return c, nil
		}
	}
}

func awsManagedClusterCreator(data *resources.TemplateData) reconciling.NamedAWSManagedClusterCreatorGetter {
	return func() (string, reconciling.AWSManagedClusterCreator) {
		return data.Bootstrap.Spec.ClusterName, func(c *awsinfrastructure.AWSManagedCluster) (*awsinfrastructure.AWSManagedCluster, error) {
			c.Name = data.Bootstrap.Spec.ClusterName
			c.Namespace = data.Namespace
			c.Spec = awsinfrastructure.AWSManagedClusterSpec{}

			return c, nil
		}
	}
}

func awsManageControlPlaneCreator(data *resources.TemplateData) reconciling.NamedAWSManagedControlPlaneCreatorGetter {
	return func() (string, reconciling.AWSManagedControlPlaneCreator) {
		return fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane"), func(c *awscontrolplane.AWSManagedControlPlane) (*awscontrolplane.AWSManagedControlPlane, error) {
			c.Name = fmt.Sprintf("%s-%s", data.Bootstrap.Spec.ClusterName, "control-plane")
			c.Namespace = data.Namespace
			c.Spec = awscontrolplane.AWSManagedControlPlaneSpec{
				EKSClusterName: data.Bootstrap.Spec.ClusterName,
				Region:         data.Bootstrap.Spec.CloudSpec.AWS.Region,
				SSHKeyName:     resources.StrPtr("default"),
				Version:        resources.StrPtr(data.Bootstrap.Spec.KubernetesVersion),
			}
			c.Spec.Addons = &[]awscontrolplane.Addon{}

			for _, addon := range data.Bootstrap.Spec.CloudSpec.AWS.Addons {
				newAddon := awscontrolplane.Addon{
					Name:                  addon.Name,
					Version:               addon.Version,
					ServiceAccountRoleArn: addon.ServiceAccountRoleArn,
				}
				*c.Spec.Addons = append(*c.Spec.Addons, newAddon)
			}
			c.Spec.AssociateOIDCProvider = true

			return c, nil
		}
	}
}

func awsMachinePoolCreator(mp bv1alpha1.AWSMachinePool, namespace, clusterName string) reconciling.NamedMachinePoolCreatorGetter {
	return func() (string, reconciling.MachinePoolCreator) {
		return mp.Name, func(c *clusterapiexp.MachinePool) (*clusterapiexp.MachinePool, error) {
			name := mp.Name
			c.Name = name
			c.Namespace = namespace
			c.Spec = clusterapiexp.MachinePoolSpec{
				ClusterName: clusterName,
				Replicas:    mp.Replicas,
				Template: clusterapi.MachineTemplateSpec{
					Spec: clusterapi.MachineSpec{
						Bootstrap: clusterapi.Bootstrap{
							DataSecretName: resources.StrPtr(""),
						},
						ClusterName: clusterName,
						InfrastructureRef: corev1.ObjectReference{
							Kind:       "AWSManagedMachinePool",
							Name:       name,
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
						},
					},
				},
			}
			return c, nil
		}
	}
}

func awsAWSManagedMachinePoolCreator(mp bv1alpha1.AWSMachinePool, namespace string) reconciling.NamedAWSManagedMachinePoolCreatorGetter {
	return func() (string, reconciling.AWSManagedMachinePoolCreator) {
		return mp.Name, func(c *awsmachinepool.AWSManagedMachinePool) (*awsmachinepool.AWSManagedMachinePool, error) {
			c.Name = mp.Name
			c.Namespace = namespace

			c.Spec = awsmachinepool.AWSManagedMachinePoolSpec{
				EKSNodegroupName:       mp.EKSNodegroupName,
				AvailabilityZones:      mp.AvailabilityZones,
				SubnetIDs:              mp.SubnetIDs,
				AdditionalTags:         infrav1.Tags(mp.AdditionalTags),
				RoleAdditionalPolicies: mp.RoleAdditionalPolicies,
				RoleName:               mp.RoleName,
				AMIVersion:             mp.AMIVersion,
				Labels:                 mp.Labels,
				DiskSize:               mp.DiskSize,
				InstanceType:           mp.InstanceType,
				ProviderIDList:         mp.ProviderIDList,
			}
			if mp.UpdateConfig != nil {
				c.Spec.UpdateConfig = &awsmachinepool.UpdateConfig{
					MaxUnavailable:           mp.UpdateConfig.MaxUnavailable,
					MaxUnavailablePercentage: mp.UpdateConfig.MaxUnavailablePercentage,
				}
			}
			if mp.CapacityType != nil {
				capacityType := awsmachinepool.ManagedMachinePoolCapacityType(*mp.CapacityType)
				c.Spec.CapacityType = &capacityType
			}
			if mp.RemoteAccess != nil {
				c.Spec.RemoteAccess = &awsmachinepool.ManagedRemoteAccess{
					SSHKeyName:           mp.RemoteAccess.SSHKeyName,
					SourceSecurityGroups: mp.RemoteAccess.SourceSecurityGroups,
					Public:               mp.RemoteAccess.Public,
				}
			}

			if mp.AMIType != nil {
				amiType := awsmachinepool.ManagedMachineAMIType(*mp.AMIType)
				c.Spec.AMIType = &amiType
			}
			if mp.Scaling != nil {
				c.Spec.Scaling = &awsmachinepool.ManagedMachinePoolScaling{
					MinSize: mp.Scaling.MinSize,
					MaxSize: mp.Scaling.MaxSize,
				}
			}

			for _, taint := range mp.Taints {
				c.Spec.Taints = append(c.Spec.Taints, awsmachinepool.Taint{
					Effect: awsmachinepool.TaintEffect(taint.Effect),
					Key:    taint.Key,
					Value:  taint.Value,
				})
			}

			return c, nil
		}
	}
}

func (aws *AWSProvider) installSA(serviceAccounts []bv1alpha1.ClusterIAMServiceAccount) error {
	os.Setenv("AWS_ACCESS_KEY_ID", aws.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", aws.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", aws.SessionToken)
	os.Setenv("AWS_REGION", aws.Region)

	aws.Data.Log.Info("Installing SA ...")
	cmd := &cmdutils.Cmd{}
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg
	cmd.ClusterConfig.Metadata.Name = aws.Data.Bootstrap.Spec.ClusterName
	cmd.ClusterConfig.Metadata.Region = aws.Region
	cmd.ProviderConfig.WaitTimeout = time.Minute * 5

	cfg.IAM.WithOIDC = api.Enabled()

	for _, sa := range serviceAccounts {
		serviceAccount := &api.ClusterIAMServiceAccount{
			ClusterIAMMeta: api.ClusterIAMMeta{
				Name:        sa.Name,
				Namespace:   sa.Namespace,
				Labels:      sa.Labels,
				Annotations: sa.Annotations,
			},
			AttachPolicyARNs: sa.AttachPolicyARNs,
			WellKnownPolicies: api.WellKnownPolicies{
				ImageBuilder:              sa.WellKnownPolicies.ImageBuilder,
				AutoScaler:                sa.WellKnownPolicies.AutoScaler,
				AWSLoadBalancerController: sa.WellKnownPolicies.AWSLoadBalancerController,
				ExternalDNS:               sa.WellKnownPolicies.ExternalDNS,
				CertManager:               sa.WellKnownPolicies.CertManager,
				EBSCSIController:          sa.WellKnownPolicies.EBSCSIController,
				EFSCSIController:          sa.WellKnownPolicies.EFSCSIController,
			},
			RoleOnly:            api.Enabled(),
			AttachRoleARN:       sa.AttachRoleARN,
			PermissionsBoundary: sa.PermissionsBoundary,
			RoleName:            sa.RoleName,
			Tags:                sa.Tags,
		}
		if !sa.RoleOnly {
			serviceAccount.RoleOnly = api.Disabled()
		}

		if sa.AttachPolicy != "" {
			var attachPolicy map[string]interface{}
			if err := yaml.Unmarshal([]byte(sa.AttachPolicy), &attachPolicy); err != nil {
				return err
			}
			serviceAccount.AttachPolicy = attachPolicy
		}

		cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)
	}

	saFilter := filter.NewIAMServiceAccountFilter()

	ctx := aws.Data.Ctx
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	oidc, err := ctl.NewOpenIDConnectManager(ctx, cfg)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		return err
	}

	if !providerExists {
		return fmt.Errorf("unable to create iamserviceaccount(s) without IAM OIDC provider enabled")
	}
	stackManager := ctl.NewStackManager(cfg)

	if err := saFilter.SetExcludeExistingFilter(ctx, stackManager, clientSet, cfg.IAM.ServiceAccounts, true); err != nil {
		return err
	}

	filteredServiceAccounts := saFilter.FilterMatching(cfg.IAM.ServiceAccounts)
	saFilter.LogInfo(cfg.IAM.ServiceAccounts)
	if filteredServiceAccounts == nil {
		existingIAMStacks, err := stackManager.ListStacksMatching(ctx, "eksctl-.*-iamserviceaccount")
		if err != nil {
			return err
		}
		return irsa.New(cfg.Metadata.Name, stackManager, oidc, clientSet).UpdateIAMServiceAccounts(ctx, cfg.IAM.ServiceAccounts, existingIAMStacks, cmd.Plan)
	}
	if err := irsa.New(cfg.Metadata.Name, stackManager, oidc, clientSet).CreateIAMServiceAccount(filteredServiceAccounts, cmd.Plan); err != nil {
		return err
	}

	return nil

}

func (aws *AWSProvider) postInstall() error {
	if len(aws.Data.Bootstrap.Spec.CloudSpec.AWS.ServiceAccounts) == 0 {
		return nil
	}
	return aws.installSA(aws.Data.Bootstrap.Spec.CloudSpec.AWS.ServiceAccounts)
}

func (aws *AWSProvider) MigrateCluster() (*ctrl.Result, error) {
	serviceAccounts := []bv1alpha1.ClusterIAMServiceAccount{
		{
			ClusterIAMMeta: bv1alpha1.ClusterIAMMeta{
				Name:      "capa-controller-manager",
				Namespace: aws.Data.Namespace,
			},
			AttachPolicyARNs:  []string{"arn:aws:iam::aws:policy/AdministratorAccess"},
			WellKnownPolicies: bv1alpha1.WellKnownPolicies{},
			RoleName:          "capa-controller-manager",
			RoleOnly:          true,
		},
	}

	if err := aws.installSA(serviceAccounts); err != nil {
		return &ctrl.Result{
			RequeueAfter: 5 * time.Second,
		}, nil
	}
	return nil, nil
}
