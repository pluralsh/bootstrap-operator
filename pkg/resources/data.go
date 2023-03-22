package resources

import (
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateData struct {
	Client    ctrlruntimeclient.Client
	Bootstrap *bv1alpha1.Bootstrap
	Namespace string
}
