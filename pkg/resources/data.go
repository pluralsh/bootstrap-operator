package resources

import (
	"context"

	"github.com/go-logr/logr"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateData struct {
	Ctx       context.Context
	Client    ctrlruntimeclient.Client
	Bootstrap *bv1alpha1.Bootstrap
	Namespace string
	Log       logr.Logger
}
