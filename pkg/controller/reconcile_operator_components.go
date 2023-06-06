package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/providers"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	initapi "github.com/pluralsh/bootstrap-operator/pkg/resources/init"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	corev1 "k8s.io/api/core/v1"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) checkOperatorComponents(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	var cp clusterapioperator.CoreProvider
	var ip clusterapioperator.InfrastructureProvider
	var bp clusterapioperator.BootstrapProvider
	var cpp clusterapioperator.ControlPlaneProvider

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

	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: resources.CoreProviderName}, &cp); err != nil {
		return nil, err
	}
	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: provider.Name()}, &ip); err != nil {
		return nil, err
	}
	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: resources.BootstrapProviderName}, &bp); err != nil {
		return nil, err
	}
	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: resources.ControlPlaneName}, &cpp); err != nil {
		return nil, err
	}

	if isReady(cp.Status.Conditions) && isReady(ip.Status.Conditions) && isReady(bp.Status.Conditions) && isReady(cpp.Status.Conditions) {
		if err := r.updateOperatorComponentsStatus(ctx, bootstrap, bv1alpha1.Running, "operator components ready", true); err != nil {
			return nil, err
		}
	}
	if err := r.checkErrors(ctx, bootstrap, []clusterv1.Conditions{cp.Status.Conditions, ip.Status.Conditions, bp.Status.Conditions, cpp.Status.Conditions}); err != nil {
		return nil, err
	}

	return &ctrl.Result{
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *Reconciler) checkErrors(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, conditions []clusterv1.Conditions) error {
	for _, condition := range conditions {
		for _, cond := range condition {
			if cond.Severity == clusterv1.ConditionSeverityError {
				if err := r.updateOperatorComponentsStatus(ctx, bootstrap, bv1alpha1.Failed, cond.Message, false); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

func isReady(conditions clusterv1.Conditions) bool {
	for _, cond := range conditions {
		if cond.Type == clusterv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (r *Reconciler) reconcileOperatorComponents(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	log := log.FromContext(ctx)
	data := &resources.TemplateData{
		Ctx:       ctx,
		Bootstrap: bootstrap,
		Namespace: r.Namespace,
		Client:    r.Client,
		Log:       log,
	}

	bootstrapProviderCreator := []reconciling.NamedBootstrapProviderCreatorGetter{
		initapi.BootstrapCreator(data),
	}
	controlplaneProviderCreator := []reconciling.NamedControlPlaneProviderCreatorGetter{
		initapi.ControlPlaneCreator(data),
	}
	infrastructureProviderCreator := []reconciling.NamedInfrastructureProviderCreatorGetter{
		initapi.InfrastructureCreator(data),
	}
	coreProviderCreator := []reconciling.NamedCoreProviderCreatorGetter{
		initapi.CoreCreator(data),
	}
	if err := reconciling.ReconcileCoreProviders(ctx, coreProviderCreator, r.Namespace, r.Client); err != nil {
		return nil, err
	}
	if err := reconciling.ReconcileInfrastructureProviders(ctx, infrastructureProviderCreator, r.Namespace, r.Client); err != nil {
		return nil, err
	}
	if err := reconciling.ReconcileBootstrapProviders(ctx, bootstrapProviderCreator, r.Namespace, r.Client); err != nil {
		return nil, err
	}
	if err := reconciling.ReconcileControlPlaneProviders(ctx, controlplaneProviderCreator, r.Namespace, r.Client); err != nil {
		return nil, err
	}

	return r.checkOperatorComponents(ctx, bootstrap)
}

func (r *Reconciler) updateOperatorComponentsStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r.Client, bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiOperatorComponentsStatus == nil {
			c.Status.CapiOperatorComponentsStatus = &bv1alpha1.Status{}
		}
		c.Status.CapiOperatorComponentsStatus.Message = message
		c.Status.CapiOperatorComponentsStatus.Phase = phase
		c.Status.CapiOperatorComponentsStatus.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update bootstrap: %w", message, err)
	}

	return nil
}
