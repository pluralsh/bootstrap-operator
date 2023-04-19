package initapi

import (
	"fmt"

	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"

	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

func CoreCreator(data *resources.TemplateData) reconciling.NamedCoreProviderCreatorGetter {
	return func() (string, reconciling.CoreProviderCreator) {
		provider, err := providers.GetProvider(data)
		if err != nil {
			return "", func(c *clusterapioperator.CoreProvider) (*clusterapioperator.CoreProvider, error) {
				return nil, err
			}
		}
		return resources.CoreProviderName, func(c *clusterapioperator.CoreProvider) (*clusterapioperator.CoreProvider, error) {
			version := data.Bootstrap.Spec.ClusterAPI.Version
			core := data.Bootstrap.Spec.ClusterAPI.Components.Core
			if core != nil && len(core.Version) > 0 {
				version = core.Version
			}

			c.Name = resources.CoreProviderName
			c.Namespace = data.Namespace
			c.Spec.SecretName = provider.Secret()
			c.Spec.Version = version

			if core != nil && len(core.FetchConfigURL) > 0 {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: fmt.Sprintf("%s/%s/core-components.yaml", core.FetchConfigURL, version),
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
