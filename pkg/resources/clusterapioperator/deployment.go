package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultResourceRequirements = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("100Mi"),
			corev1.ResourceCPU:    resource.MustParse("100m"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("150Mi"),
			corev1.ResourceCPU:    resource.MustParse("150m"),
		},
	}
)

func DeploymentCreator(data *resources.TemplateData) reconciling.NamedDeploymentCreatorGetter {
	return func() (string, reconciling.DeploymentCreator) {
		return resources.ClusterAPIOperatorDeploymentName, func(dep *appsv1.Deployment) (*appsv1.Deployment, error) {
			additionalLabels := map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator", "control-plane": "controller-manager"}

			dep.Name = resources.ClusterAPIOperatorDeploymentName
			dep.Labels = resources.BaseAppLabels(resources.ClusterAPIOperatorDeploymentName, additionalLabels)

			dep.Spec.Replicas = resources.Int32(1)

			dep.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: resources.BaseAppLabels(resources.ClusterAPIOperatorDeploymentName, additionalLabels),
			}

			dep.Spec.Template = corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: resources.BaseAppLabels(resources.ClusterAPIOperatorDeploymentName, additionalLabels),
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: resources.Int64(10),
					Tolerations: []corev1.Toleration{
						{
							Key:    "node-role.kubernetes.io/master",
							Effect: "NoSchedule",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "cert",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  resources.CertManagerWebhookSecretName,
									DefaultMode: resources.Int32(420),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "manager",
							Image:   data.Bootstrap.Spec.ClusterAPI.Components.Operator.ManagerImage,
							Command: []string{"/manager"},
							Args:    []string{"--metrics-bind-addr", "127.0.0.1:8080", "--leader-elect"},
							Ports: []corev1.ContainerPort{
								{
									Name:          "webhook-server",
									ContainerPort: 9443,
									Protocol:      "TCP",
								},
							},
							Resources: defaultResourceRequirements,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "cert",
									ReadOnly:  true,
									MountPath: "/tmp/k8s-webhook-server/serving-certs",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						{
							Name:  "kube-rbac-proxy",
							Image: data.Bootstrap.Spec.ClusterAPI.Components.Operator.KubeRBACProxyImage,
							Args:  []string{"--secure-listen-address", "0.0.0.0:8443", "--upstream", "http://127.0.0.1:8080/", "--logtostderr", "true", "--v", "10"},
							Ports: []corev1.ContainerPort{
								{
									Name:          "https",
									ContainerPort: 8443,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: resources.Disabled(),
								RunAsUser:                resources.Int64(65532),
							},
						},
					},
				},
			}

			return dep, nil
		}
	}
}
