package providers

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/pluralsh/bootstrap-operator/pkg/resources"
)

type Provider interface {
	Name() string
	Secret() string
	Version() string
	FetchConfigURL() string
	ReconcileCluster() error
	CheckCluster() (*ctrl.Result, error)
	Init() (*ctrl.Result, error)
}

func GetProvider(data *resources.TemplateData) (Provider, error) {
	if data.Bootstrap.Spec.CloudSpec.AWS != nil {
		return GetAWSProvider(data)
	}
	if data.Bootstrap.Spec.CloudSpec.GCP != nil {
		return GetGCPProvider(data)
	}
	return nil, fmt.Errorf("invalid provider")
}
