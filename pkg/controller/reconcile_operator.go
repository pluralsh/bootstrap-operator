package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/clusterapioperator"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) reconcileOperator(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) (*ctrl.Result, error) {
	if err := r.ensureOperatorResourcesAreDeployed(ctx, bootstrap, r.Namespace); err != nil {
		return nil, err
	}

	var deployment appsv1.Deployment
	if err := r.Get(ctx, ctrlruntimeclient.ObjectKey{Namespace: r.Namespace, Name: resources.ClusterAPIOperatorDeploymentName}, &deployment); err != nil {
		return nil, err
	}
	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		if err := r.updateOperatorStatus(ctx, bootstrap, bv1alpha1.Running, "cluster API operator is up and running", true); err != nil {
			return nil, err
		}
		return nil, nil
	}

	return &ctrl.Result{
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *Reconciler) ensureOperatorResourcesAreDeployed(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, namespace string) error {
	data := &resources.TemplateData{
		Bootstrap: bootstrap,
		Namespace: namespace,
	}

	serviceCreators := []reconciling.NamedServiceCreatorGetter{
		clusterapioperator.ServiceCreator(),
		clusterapioperator.WebhookServiceCreator(),
	}
	roleCreators := []reconciling.NamedRoleCreatorGetter{
		clusterapioperator.LeaderElectionRoleCreator(data),
	}
	clusterRoleCreators := []reconciling.NamedClusterRoleCreatorGetter{
		clusterapioperator.ManagerClusterRoleCreator(),
		clusterapioperator.MetricsClusterRoleCreator(),
		clusterapioperator.ProxyClusterRoleCreator(),
	}
	roleBindingCreators := []reconciling.NamedRoleBindingCreatorGetter{
		clusterapioperator.LeaderElectionRoleBindingCreator(data),
	}
	clusterRoleBindingCreators := []reconciling.NamedClusterRoleBindingCreatorGetter{
		clusterapioperator.ManagerClusterRoleBindingCreator(data),
		clusterapioperator.ProxyClusterRoleBindingCreator(data),
	}
	deploymentCreators := []reconciling.NamedDeploymentCreatorGetter{
		clusterapioperator.DeploymentCreator(data),
	}
	certCreators := []reconciling.NamedCertificateCreatorGetter{
		clusterapioperator.CertificateCreator(data),
	}
	issuerCreators := []reconciling.NamedIssuerCreatorGetter{
		clusterapioperator.IssuerCreator(data),
	}
	validatingWebhookConfigurationCreators := []reconciling.NamedValidatingWebhookConfigurationCreatorGetter{
		clusterapioperator.ValidatingWebhookConfigurationCreator(data),
	}
	if err := reconciling.ReconcileServices(ctx, serviceCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileRoles(ctx, roleCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileClusterRoles(ctx, clusterRoleCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileRoleBindings(ctx, roleBindingCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileClusterRoleBindings(ctx, clusterRoleBindingCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileDeployments(ctx, deploymentCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileCertificates(ctx, certCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileIssuers(ctx, issuerCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	if err := reconciling.ReconcileValidatingWebhookConfigurations(ctx, validatingWebhookConfigurationCreators, r.Namespace, r.Client); err != nil {
		return err
	}
	return nil
}

func (r *Reconciler) updateOperatorStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r.Client, bootstrap, func(c *bv1alpha1.Bootstrap) {
		if c.Status.CapiOperatorStatus == nil {
			c.Status.CapiOperatorStatus = &bv1alpha1.Status{}
		}
		c.Status.CapiOperatorStatus.Message = message
		c.Status.CapiOperatorStatus.Phase = phase
		c.Status.CapiOperatorStatus.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update cluster: %w", message, err)
	}

	return nil
}
