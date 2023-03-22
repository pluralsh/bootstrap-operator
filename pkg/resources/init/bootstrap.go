package initapi

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

func BootstrapCreator(data *resources.TemplateData) reconciling.NamedBootstrapProviderCreatorGetter {
	return func() (string, reconciling.BootstrapProviderCreator) {
		return resources.BootstrapProviderName, func(c *clusterapioperator.BootstrapProvider) (*clusterapioperator.BootstrapProvider, error) {
			c.Name = resources.BootstrapProviderName
			c.Namespace = data.Namespace
			if data.Bootstrap.Spec.Components.Bootstrap.Version != "" {
				c.Spec.Version = data.Bootstrap.Spec.Components.Bootstrap.Version
			}
			if data.Bootstrap.Spec.Components.Bootstrap.FetchConfigURL != "" {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: data.Bootstrap.Spec.Components.Bootstrap.FetchConfigURL,
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
