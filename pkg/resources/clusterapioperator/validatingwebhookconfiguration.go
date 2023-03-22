package clusterapioperator

import (
	"fmt"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

func ValidatingWebhookConfigurationCreator(data *resources.TemplateData) reconciling.NamedValidatingWebhookConfigurationCreatorGetter {
	return func() (string, reconciling.ValidatingWebhookConfigurationCreator) {
		return resources.ClusterAPIValidatingWebhookConfigurationName, func(i *admissionregistrationv1.ValidatingWebhookConfiguration) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
			failurePolicy := admissionregistrationv1.Fail
			sideEffect := admissionregistrationv1.SideEffectClassNone
			i.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			i.Name = resources.ClusterAPIValidatingWebhookConfigurationName

			i.Annotations = map[string]string{"cert-manager.io/inject-ca-from": fmt.Sprintf("%s/capi-operator-serving-cert", data.Namespace)}
			i.Webhooks = []admissionregistrationv1.ValidatingWebhook{
				{
					Name: "vbootstrapprovider.kb.io",
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Namespace: data.Namespace,
							Name:      resources.ClusterAPIOperatorWebhookServiceName,
							Path:      resources.StrPtr("/validate-operator-cluster-x-k8s-io-v1alpha1-bootstrapprovider"),
						},
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Operations: []admissionregistrationv1.OperationType{
								admissionregistrationv1.Create,
								admissionregistrationv1.Update,
							},
							Rule: admissionregistrationv1.Rule{
								APIGroups:   []string{"operator.cluster.x-k8s.io"},
								APIVersions: []string{"v1alpha1"},
								Resources:   []string{"bootstrapproviders"},
							},
						},
					},
					SideEffects:   &sideEffect,
					FailurePolicy: &failurePolicy,
					AdmissionReviewVersions: []string{
						"v1",
						"v1alpha1",
					},
				},
				{
					Name: "vcontrolplaneprovider.kb.io",
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Namespace: data.Namespace,
							Name:      resources.ClusterAPIOperatorWebhookServiceName,
							Path:      resources.StrPtr("/validate-operator-cluster-x-k8s-io-v1alpha1-controlplaneprovider"),
						},
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Operations: []admissionregistrationv1.OperationType{
								admissionregistrationv1.Create,
								admissionregistrationv1.Update,
							},
							Rule: admissionregistrationv1.Rule{
								APIGroups:   []string{"operator.cluster.x-k8s.io"},
								APIVersions: []string{"v1alpha1"},
								Resources:   []string{"controlplaneproviders"},
							},
						},
					},
					SideEffects:   &sideEffect,
					FailurePolicy: &failurePolicy,
					AdmissionReviewVersions: []string{
						"v1",
						"v1alpha1",
					},
				},
				{
					Name: "vcoreprovider.kb.io",
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Namespace: data.Namespace,
							Name:      resources.ClusterAPIOperatorWebhookServiceName,
							Path:      resources.StrPtr("/validate-operator-cluster-x-k8s-io-v1alpha1-coreprovider"),
						},
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Operations: []admissionregistrationv1.OperationType{
								admissionregistrationv1.Create,
								admissionregistrationv1.Update,
							},
							Rule: admissionregistrationv1.Rule{
								APIGroups:   []string{"operator.cluster.x-k8s.io"},
								APIVersions: []string{"v1alpha1"},
								Resources:   []string{"coreproviders"},
							},
						},
					},
					SideEffects:   &sideEffect,
					FailurePolicy: &failurePolicy,
					AdmissionReviewVersions: []string{
						"v1",
						"v1alpha1",
					},
				},
				{
					Name: "vinfrastructureprovider.kb.io",
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Namespace: data.Namespace,
							Name:      resources.ClusterAPIOperatorWebhookServiceName,
							Path:      resources.StrPtr("/validate-operator-cluster-x-k8s-io-v1alpha1-infrastructureprovider"),
						},
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Operations: []admissionregistrationv1.OperationType{
								admissionregistrationv1.Create,
								admissionregistrationv1.Update,
							},
							Rule: admissionregistrationv1.Rule{
								APIGroups:   []string{"operator.cluster.x-k8s.io"},
								APIVersions: []string{"v1alpha1"},
								Resources:   []string{"infrastructureproviders"},
							},
						},
					},
					SideEffects:   &sideEffect,
					FailurePolicy: &failurePolicy,
					AdmissionReviewVersions: []string{
						"v1",
						"v1alpha1",
					},
				},
			}
			return i, nil
		}
	}
}
