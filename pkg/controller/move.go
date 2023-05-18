package controller

import (
	"context"
	"fmt"

	coreapi "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	awsinfrastructure "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	awscontrolplane "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	awsmachinepool "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	azurecontroleplane "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	gcpclusterapi "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterapi "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/move"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
)

func (r *Reconciler) moveNamespace(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) error {
	if err := r.move(ctx, bootstrap); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) move(ctx context.Context, bootstrap *bv1alpha1.Bootstrap) error {
	r.Log.Info("Move CAPI components")
	var kc coreapi.Secret
	if err := r.Client.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-kubeconfig", bootstrap.Spec.ClusterName), Namespace: r.Namespace}, &kc); err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("CAPI components already moved")
			return nil
		}
		return err
	}

	namespaceCreator := []reconciling.NamedNamespaceCreatorGetter{
		namespaceCreator(r.Namespace),
	}

	newCluster := move.Cluster{
		Namespace:   r.Namespace,
		ClusterName: bootstrap.Spec.ClusterName,
		Scheme:      r.Scheme,
		Client:      r.Client,
		Ctx:         ctx,
		Log:         r.Log,
		Bootstrap:   bootstrap,
	}

	toClient, err := newCluster.GetClient()
	if err != nil {
		return err
	}
	r.Log.Info("toClient created successfully")

	if err := reconciling.ReconcileNamespaces(ctx, namespaceCreator, r.Namespace, toClient); err != nil {
		return err
	}
	crdList := &apiextensionsv1.CustomResourceDefinitionList{}
	if err := getCRDList(ctx, r.Client, crdList); err != nil {
		return err
	}

	allowedGroups := sets.NewString(clusterapi.GroupVersion.Group,
		awsinfrastructure.GroupVersion.Group,
		awscontrolplane.GroupVersion.Group,
		awsmachinepool.GroupVersion.Group,
		azurecontroleplane.GroupVersion.Group,
		"aadpodidentity.k8s.io", // TODO: Remove it once CAPI will be the default.
		clusterapiexp.GroupVersion.Group,
		gcpclusterapi.GroupVersion.Group)
	var crdCreatorGetter []reconciling.NamedCustomResourceDefinitionCreatorGetter
	for _, crd := range crdList.Items {
		if allowedGroups.Has(crd.Spec.Group) {
			crdCreatorGetter = append(crdCreatorGetter, crdCreator(crd))
		}
	}
	if err := reconciling.ReconcileCustomResourceDefinitions(ctx, crdCreatorGetter, r.Namespace, toClient); err != nil {
		return err
	}

	if err := newCluster.MoveClusterAPI(); err != nil {
		return err
	}
	return nil
}

func getCRDList(ctx context.Context, client client.Client, crdList *apiextensionsv1.CustomResourceDefinitionList) error {
	if err := client.List(ctx, crdList); err != nil {
		return err
	}
	return nil
}

func crdCreator(crd apiextensionsv1.CustomResourceDefinition) reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return crd.Name, func(c *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			c = crd.DeepCopy()
			c.ObjectMeta.ResourceVersion = ""
			return c, nil
		}
	}
}

func namespaceCreator(namespace string) reconciling.NamedNamespaceCreatorGetter {
	return func() (string, reconciling.NamespaceCreator) {
		return namespace, func(c *coreapi.Namespace) (*coreapi.Namespace, error) {
			c.Name = namespace
			return c, nil
		}
	}
}
