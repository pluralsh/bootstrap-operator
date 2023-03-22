package clusterapioperator

import (
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

func IssuerCreator(data *resources.TemplateData) reconciling.NamedIssuerCreatorGetter {
	return func() (string, reconciling.IssuerCreator) {
		return resources.ClusterAPIIssuerName, func(i *v1.Issuer) (*v1.Issuer, error) {
			i.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			i.Name = resources.ClusterAPIIssuerName
			i.Namespace = data.Namespace
			i.Spec.SelfSigned = &v1.SelfSignedIssuer{}
			return i, nil
		}
	}
}
