package controller

import (
	"context"
	"fmt"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) reconcileCluster(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	provider, err := providers.GetProvider(&resources.TemplateData{
		Ctx:       ctx,
		Client:    r.Client,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Log:       r.Log,
	})
	if err != nil {
		return nil, err
	}
	if err := provider.ReconcileCluster(); err != nil {
		return nil, err
	}
	return provider.CheckCluster()
}

func (r *Reconciler) checkCluster(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	provider, err := providers.GetProvider(&resources.TemplateData{
		Ctx:       ctx,
		Client:    r.Client,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Log:       r.Log,
	})
	if err != nil {
		return nil, err
	}
	return provider.CheckCluster()
}

func (r *Reconciler) updateClusterStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r.Client, bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiClusterStatus == nil {
			c.Status.CapiClusterStatus = &bv1alpha1.ClusterStatus{}
		}
		c.Status.CapiClusterStatus.Message = message
		c.Status.CapiClusterStatus.Phase = phase
		c.Status.CapiClusterStatus.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update bootstrap: %w", message, err)
	}

	return nil
}
