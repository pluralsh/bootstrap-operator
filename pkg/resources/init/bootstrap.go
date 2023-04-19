package initapi

import (
	"fmt"

	"github.com/pluralsh/bootstrap-operator/pkg/providers"

	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"

	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
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
			version := data.Bootstrap.Spec.ClusterAPI.Version
			bootstrap := data.Bootstrap.Spec.ClusterAPI.Components.Bootstrap
			if bootstrap != nil && len(bootstrap.Version) > 0 {
				version = bootstrap.Version
			}

			c.Name = resources.BootstrapProviderName
			c.Namespace = data.Namespace
			c.Spec.SecretName = provider.Secret()
			c.Spec.Version = version

			if bootstrap != nil && len(bootstrap.FetchConfigURL) > 0 {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: fmt.Sprintf("%s/%s/bootstrap-components.yaml", bootstrap.FetchConfigURL, version),
				}
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
