package clusterapioperator

import (
	"fmt"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

func CertificateCreator(data *resources.TemplateData) reconciling.NamedCertificateCreatorGetter {
	return func() (string, reconciling.CertificateCreator) {
		return resources.ClusterAPICertificateName, func(c *v1.Certificate) (*v1.Certificate, error) {
			c.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			c.Name = resources.ClusterAPICertificateName
			c.Namespace = data.Namespace
			c.Spec = v1.CertificateSpec{
				DNSNames: []string{
					fmt.Sprintf("capi-operator-webhook-service.%s.svc", data.Namespace),
					fmt.Sprintf("capi-operator-webhook-service.%s.svc.cluster.local", data.Namespace),
				},
				IssuerRef: cmmeta.ObjectReference{
					Name: "capi-operator-selfsigned-issuer",
					Kind: "Issuer",
				},
				SecretName: resources.CertManagerWebhookSecretName,
			}
			return c, nil
		}
	}
}
