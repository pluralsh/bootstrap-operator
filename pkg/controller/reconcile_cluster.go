package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) reconcileCluster(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	provider, err := providers.GetProvider(&resources.TemplateData{
		Ctx:       ctx,
		Client:    r.Client,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Log:       log,
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
	log := log.FromContext(ctx)
	provider, err := providers.GetProvider(&resources.TemplateData{
		Ctx:       ctx,
		Client:    r.Client,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Log:       log,
	})
	if err != nil {
		return nil, err
	}
	return provider.CheckCluster()
}

func (r *Reconciler) migration(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	provider, err := providers.GetProvider(&resources.TemplateData{
		Ctx:       ctx,
		Client:    r.Client,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Log:       log,
	})
	if err != nil {
		return nil, err
	}
	pods := &corev1.PodList{}
	selector := fmt.Sprintf("infrastructure-%s", provider.Name())
	if err := r.List(context.Background(), pods, ctrlruntimeclient.MatchingLabels{"cluster.x-k8s.io/provider": selector}); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get pods: %w", err)
		}
		log.Info("Waiting for infrastructure operator ...")
		return &ctrl.Result{
			RequeueAfter: 5 * time.Second,
		}, nil
	}
	if len(pods.Items) > 0 {
		if isPodReady(pods.Items[0].Status.Conditions) {
			log.Info("Infrastructure operator ready")
			return provider.MigrateCluster()
		}
	}
	return &ctrl.Result{
		RequeueAfter: 5 * time.Second,
	}, nil
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

func isPodReady(conditions []corev1.PodCondition) bool {
	for _, cond := range conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
