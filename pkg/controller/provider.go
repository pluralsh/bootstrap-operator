package controller

import (
	"context"
	"fmt"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) initProvider(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	data := &resources.TemplateData{
		Ctx:       ctx,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Client:    r.Client,
		Log:       log,
	}
	prov, err := providers.GetProvider(data)
	if err != nil {
		return nil, err
	}
	return prov.Init()
}

func (r *Reconciler) updateProviderStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r.Client, bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.ProviderStatus == nil {
			c.Status.ProviderStatus = &bv1alpha1.Status{}
		}
		c.Status.ProviderStatus.Message = message
		c.Status.ProviderStatus.Phase = phase
		c.Status.ProviderStatus.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update bootstrap: %w", message, err)
	}

	return nil
}
