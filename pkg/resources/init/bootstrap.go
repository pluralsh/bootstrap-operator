package initapi

import (
	"fmt"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"

	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

func BootstrapCreator(data *resources.TemplateData) reconciling.NamedBootstrapProviderCreatorGetter {
	return func() (string, reconciling.BootstrapProviderCreator) {
		provider, err := providers.GetProvider(data)
		if err != nil {
			return "", func(c *clusterapioperator.BootstrapProvider) (*clusterapioperator.BootstrapProvider, error) {
				return nil, err
			}
		}
		return resources.BootstrapProviderName, func(c *clusterapioperator.BootstrapProvider) (*clusterapioperator.BootstrapProvider, error) {
			c.Name = resources.BootstrapProviderName
			c.Namespace = data.Namespace
			c.Spec.SecretName = provider.Secret()
			c.Spec.Version = data.Bootstrap.Spec.ClusterAPI.Version
			c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
				URL: fmt.Sprintf("https://github.com/kubernetes-sigs/cluster-api/releases/%s/bootstrap-components.yaml", data.Bootstrap.Spec.ClusterAPI.Version),
			}
			c.Spec.Deployment = &clusterapioperator.DeploymentSpec{
				Containers: []clusterapioperator.ContainerSpec{
					{
						Name: "manager",
					},
				},
			}

			return c, nil
		}
	}
}
