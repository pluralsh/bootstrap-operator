package resources

const (
	ClusterAPIOperatorDeploymentName             = "capi-operator-controller-manager"
	ClusterAPIOperatorServiceName                = "capi-operator-controller-manager-metrics-service"
	ClusterAPIOperatorWebhookServiceName         = "capi-operator-webhook-service"
	ClusterAPILeaderElectionRoleName             = "capi-operator-leader-election-role"
	ClusterAPIManagerClusterRoleName             = "capi-operator-manager-role"
	ClusterAPIMetricsClusterRoleName             = "capi-operator-metrics-reader"
	ClusterAPIProxyClusterRoleName               = "capi-operator-proxy-role"
	ClusterAPILeaderElectionRoleBindingName      = "capi-operator-leader-election-rolebinding"
	ClusterAPIManagerClusterRoleBindingName      = "capi-operator-manager-rolebinding"
	ClusterAPIProxyClusterRoleBindingName        = "capi-operator-proxy-rolebinding"
	ClusterAPICertificateName                    = "capi-operator-serving-cert"
	ClusterAPIIssuerName                         = "capi-operator-selfsigned-issuer"
	ClusterAPIValidatingWebhookConfigurationName = "capi-operator-validating-webhook-configuration"

	BootstrapProviderName = "bootstrap"
	CoreProviderName      = "cluster-api"
	ControlPlaneName      = "control-plane"

	CertManagerWebhookSecretName = "capi-operator-webhook-service-cert"

	AppLabelKey = "app"
)

// BaseAppLabels returns the minimum required labels.
func BaseAppLabels(name string, additionalLabels map[string]string) map[string]string {
	labels := map[string]string{
		AppLabelKey: name,
	}
	for k, v := range additionalLabels {
		labels[k] = v
	}
	return labels
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

func StrPtr(s string) *string {
	return &s
}
