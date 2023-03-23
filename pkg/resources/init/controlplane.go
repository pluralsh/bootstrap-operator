package initapi

import (
	"fmt"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
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
			c.Name = resources.ControlPlaneName
			c.Namespace = data.Namespace
			c.Spec.SecretName = provider.Secret()
			c.Spec.Version = data.Bootstrap.Spec.ClusterAPI.Version
			c.Spec.FetchConfig = &clusterapioperator.FetchConfiguration{
				URL: fmt.Sprintf("https://github.com/kubernetes-sigs/cluster-api/releases/%s/control-plane-components.yaml", data.Bootstrap.Spec.ClusterAPI.Version),
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
