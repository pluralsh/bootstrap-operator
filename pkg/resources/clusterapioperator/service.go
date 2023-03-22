package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ServiceCreator() reconciling.NamedServiceCreatorGetter {
	return func() (string, reconciling.ServiceCreator) {
		return resources.ClusterAPIOperatorServiceName, func(se *corev1.Service) (*corev1.Service, error) {
			labels := map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator", "control-plane": "controller-manager"}
			se.Name = resources.ClusterAPIOperatorServiceName
			se.Labels = labels
			se.Spec.Ports = []corev1.ServicePort{
				{
					Name:       "https",
					Port:       8443,
					TargetPort: intstr.FromString("https"),
				},
			}
			se.Spec.Selector = labels

			return se, nil
		}

	}

}

func WebhookServiceCreator() reconciling.NamedServiceCreatorGetter {
	return func() (string, reconciling.ServiceCreator) {
		return resources.ClusterAPIOperatorWebhookServiceName, func(se *corev1.Service) (*corev1.Service, error) {
			labels := map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator", "control-plane": "controller-manager"}
			se.Name = resources.ClusterAPIOperatorWebhookServiceName
			se.Labels = labels
			se.Spec.Ports = []corev1.ServicePort{
				{
					Port:       443,
					TargetPort: intstr.FromInt(9443),
				},
			}
			se.Spec.Selector = labels

			return se, nil
		}

	}

}
