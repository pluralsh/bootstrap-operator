package resources

import (
	"context"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"go.uber.org/zap"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateData struct {
	Ctx       context.Context
	Client    ctrlruntimeclient.Client
	Bootstrap *bv1alpha1.Bootstrap
	Namespace string
	Log       *zap.SugaredLogger
}
