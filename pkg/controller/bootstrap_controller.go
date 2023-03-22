package controller

import (
	"context"
	"fmt"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler reconciles a DatabaseRequest object
type Reconciler struct {
	client.Client
	Log       *zap.SugaredLogger
	Namespace string
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	bootstrap := &bv1alpha1.Bootstrap{}
	if err := r.Get(ctx, req.NamespacedName, bootstrap); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !bootstrap.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	res, err := r.reconcileOperator(ctx, bootstrap)
	if err != nil {
		updateErr := r.updateOperatorStatus(ctx, bootstrap, bv1alpha1.Failed, err.Error(), false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI operator: %w", err)
	}
	if res != nil {
		return *res, nil
	}
	if !bootstrap.Status.CapiOperatorStatus.Ready {
		return ctrl.Result{}, nil
	}

	res, err = r.reconcileCore(ctx, bootstrap)
	if err != nil {
		updateErr := r.updateCoreStatus(ctx, bootstrap, bv1alpha1.Failed, err.Error(), false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI core: %w", err)
	}
	if res != nil {
		return *res, nil
	}
	if !bootstrap.Status.CapiCore.Ready {
		return ctrl.Result{}, nil
	}

	res, err = r.reconcileBootstrap(ctx, bootstrap)
	if err != nil {
		updateErr := r.updateBootstrapStatus(ctx, bootstrap, bv1alpha1.Failed, err.Error(), false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI bootstrap: %w", err)
	}
	if res != nil {
		return *res, nil
	}
	if !bootstrap.Status.CapiBootstrap.Ready {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bv1alpha1.Bootstrap{}).
		Owns(&appsv1.Deployment{}).
		//Owns(&clusterapioperator.BootstrapProvider{}).
		Complete(r)
}
