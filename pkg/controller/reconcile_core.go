package controller

import (
	"context"
	"fmt"
	"time"

	initapi "github.com/pluralsh/bootstrap-operator/pkg/resources/init"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) reconcileCore(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	if bootstrap.Status.CapiCore == nil {
		if err := r.updateBootstrapStatus(ctx, bootstrap, bv1alpha1.Creating, "creating cluster API core component", false); err != nil {
			return nil, err
		}
	}
	data := &resources.TemplateData{
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
	}

	coreProviderCreator := []reconciling.NamedCoreProviderCreatorGetter{
		initapi.CoreCreator(data),
	}
	if err := reconciling.ReconcileCoreProviders(ctx, coreProviderCreator, r.Namespace, r); err != nil {
		return nil, err
	}

	var cp clusterapioperator.CoreProvider
	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: resources.CoreProviderName}, &cp); err != nil {
		return nil, err
	}

	for _, cond := range cp.Status.Conditions {
		if cond.Type == clusterv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			if err := r.updateCoreStatus(ctx, bootstrap, bv1alpha1.Running, "cluster API core is up and running", true); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}

	return &ctrl.Result{
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *Reconciler) updateCoreStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r, bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiCore == nil {
			c.Status.CapiCore = &bv1alpha1.Status{}
		}
		c.Status.CapiCore.Message = message
		c.Status.CapiCore.Phase = phase
		c.Status.CapiCore.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update cluster: %w", message, err)
	}

	return nil
}
