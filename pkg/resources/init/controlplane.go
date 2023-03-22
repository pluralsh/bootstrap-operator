package initapi

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

func ControlPlaneCreator(data *resources.TemplateData) reconciling.NamedControlPlaneProviderCreatorGetter {
	return func() (string, reconciling.ControlPlaneProviderCreator) {
		return resources.ControlPlaneName, func(c *clusterapioperator.ControlPlaneProvider) (*clusterapioperator.ControlPlaneProvider, error) {
			c.Name = resources.ControlPlaneName
			c.Namespace = data.Namespace
			if data.Bootstrap.Spec.Components.ControlPlane.Version != "" {
				c.Spec.Version = data.Bootstrap.Spec.Components.ControlPlane.Version
			}
			if data.Bootstrap.Spec.Components.ControlPlane.FetchConfigURL != "" {
				c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
					URL: data.Bootstrap.Spec.Components.ControlPlane.FetchConfigURL,
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
