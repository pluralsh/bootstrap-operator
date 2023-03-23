package providers

import (
	"fmt"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Provider interface {
	Name() string
	Secret() string
	Version() string
	ReconcileCluster() error
	CheckCluster() (*ctrl.Result, error)
	Init() (*ctrl.Result, error)
}

func GetProvider(data *resources.TemplateData) (Provider, error) {
	if data.Bootstrap.Spec.CloudSpec.AWS != nil {
		return GetAWSProvider(data)
	}
	return nil, fmt.Errorf("invalid provider")
}
