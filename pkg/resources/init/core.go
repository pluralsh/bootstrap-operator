package initapi

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

func CoreCreator(data *resources.TemplateData) reconciling.NamedCoreProviderCreatorGetter {
	return func() (string, reconciling.CoreProviderCreator) {
		return resources.CoreProviderName, func(c *clusterapioperator.CoreProvider) (*clusterapioperator.CoreProvider, error) {
			c.Name = resources.CoreProviderName
			c.Namespace = data.Namespace
			if data.Bootstrap.Spec.Components.Core.Version != "" {
				c.Spec.Version = data.Bootstrap.Spec.Components.Core.Version
			}
			if data.Bootstrap.Spec.Components.Core.FetchConfigURL != "" {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: data.Bootstrap.Spec.Components.Core.FetchConfigURL,
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
