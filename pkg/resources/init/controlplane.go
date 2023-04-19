package initapi

import (
	"fmt"

	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"

	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

func ControlPlaneCreator(data *resources.TemplateData) reconciling.NamedControlPlaneProviderCreatorGetter {
	return func() (string, reconciling.ControlPlaneProviderCreator) {
		provider, err := providers.GetProvider(data)
		if err != nil {
			return "", func(c *clusterapioperator.ControlPlaneProvider) (*clusterapioperator.ControlPlaneProvider, error) {
				return nil, err
			}
		}
		return resources.ControlPlaneName, func(c *clusterapioperator.ControlPlaneProvider) (*clusterapioperator.ControlPlaneProvider, error) {
			version := data.Bootstrap.Spec.ClusterAPI.Version
			controlPlane := data.Bootstrap.Spec.ClusterAPI.Components.ControlPlane
			if controlPlane != nil && len(controlPlane.Version) > 0 {
				version = controlPlane.Version
			}

			c.Name = resources.ControlPlaneName
			c.Namespace = data.Namespace
			c.Spec.SecretName = provider.Secret()
			c.Spec.Version = version

			if controlPlane != nil && len(controlPlane.FetchConfigURL) > 0 {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: fmt.Sprintf("%s/%s/control-plane-components.yaml", controlPlane.FetchConfigURL, version),
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
