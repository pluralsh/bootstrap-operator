package initapi

import (
	"time"

	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
	ctrlconfigv1 "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

func InfrastructureCreator(data *resources.TemplateData) reconciling.NamedInfrastructureProviderCreatorGetter {
	return func() (string, reconciling.InfrastructureProviderCreator) {
		provider, err := providers.GetProvider(data)
		if err != nil {
			return "", func(c *clusterapioperator.InfrastructureProvider) (*clusterapioperator.InfrastructureProvider, error) {
				return nil, err
			}
		}
		return provider.Name(), func(c *clusterapioperator.InfrastructureProvider) (*clusterapioperator.InfrastructureProvider, error) {
			c.Name = provider.Name()
			c.Namespace = data.Namespace
			c.Spec.Version = provider.Version()
			c.Spec.SecretName = provider.Secret()
			c.Spec.Manager = &clusterapioperator.ManagerSpec{
				ControllerManagerConfigurationSpec: ctrlconfigv1.ControllerManagerConfigurationSpec{
					SyncPeriod: &metav1.Duration{Duration: 30 * time.Second},
				},
			}
			c.Spec.Deployment = &clusterapioperator.DeploymentSpec{
				Containers: []clusterapioperator.ContainerSpec{
					{
						Name: "manager",
						Args: map[string]string{
							"awscluster-concurrency": "12",
							"awsmachine-concurrency": "11",
						},
					},
				},
			}

			return c, nil
		}
	}
}
