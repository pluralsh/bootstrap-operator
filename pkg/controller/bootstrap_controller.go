package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/bootstrap-operator/apis/bootstrap/helper"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	clusterNameField = ".spec.clusterName"
	// eksLabel         = "eks.amazonaws.com/capacityType"
	// gkeLabel         = "cloud.google.com/gke-os-distribution"
	// aksLabel         = "kubernetes.azure.com/os"
)

// Reconciler reconciles a DatabaseRequest object
type Reconciler struct {
	client.Client
	KubeClient *kubernetes.Clientset
	Namespace  string
	Scheme     *runtime.Scheme
	Kubeconfig string
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	bootstrap := &bv1alpha1.Bootstrap{}
	if err := r.Get(ctx, req.NamespacedName, bootstrap); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var kindNet appsv1.DaemonSet

	err := r.Client.Get(ctx, types.NamespacedName{Name: "kindnet", Namespace: "kube-system"}, &kindNet)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Not running in kind")
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to get kindnet daemonset")
		return ctrl.Result{}, err
	} else {
		log.Info("Running in kind")
	}

	if !bootstrap.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	log.Info("Reconciling bootstrap", "name", bootstrap.GetName())

	if !r.CheckCertManager(ctx) {
		log.Info("Waiting for cert manager")
		return ctrl.Result{
			RequeueAfter: 5 * time.Second,
		}, nil
	}

	if bootstrap.Status.CapiOperatorStatus == nil {
		if err := r.updateOperatorStatus(ctx, bootstrap, bv1alpha1.Creating, "creating CAPI operator", false); err != nil {
			log.Error(err, "failed to set the bootstrap error")
			return ctrl.Result{}, err
		}
	}
	if bootstrap.Status.CapiOperatorComponentsStatus == nil {
		if err := r.updateOperatorComponentsStatus(ctx, bootstrap, bv1alpha1.Creating, "creating CAPI operator components", false); err != nil {
			return ctrl.Result{}, err
		}
	}
	if bootstrap.Status.CapiClusterStatus == nil {
		if err := r.updateClusterStatus(ctx, bootstrap, bv1alpha1.Creating, "creating CAPI cluster", false); err != nil {
			return ctrl.Result{}, err
		}
	}

	if bootstrap.Status.ProviderStatus == nil {
		if err := r.updateProviderStatus(ctx, bootstrap, bv1alpha1.Started, "init cloud provider", false); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !bootstrap.Status.ProviderStatus.Ready {
		res, err := r.initProvider(ctx, bootstrap)
		if err != nil {
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			return ctrl.Result{}, fmt.Errorf("failed to init provider: %w", err)
		}
		updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Creating, "Init cloud operator", false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		if res != nil {
			return *res, nil
		}
	}

	if !bootstrap.Status.CapiOperatorStatus.Ready {
		res, err := r.reconcileOperator(ctx, bootstrap)
		if err != nil {
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI operator: %w", err)
		}
		updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Creating, "Creating CAPI operator", false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		if res != nil {
			return *res, nil
		}
	}

	if !bootstrap.Status.CapiOperatorComponentsStatus.Ready {
		res, err := r.reconcileOperatorComponents(ctx, bootstrap)
		if err != nil {
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI operator: %w", err)
		}
		updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Creating, "Creating CAPI operator components", false)
		if updateErr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
		}
		if res != nil {
			return *res, nil
		}
	}
	if bootstrap.Spec.MigrateCluster {
		if !bootstrap.Status.Ready {
			res, err := r.migration(ctx, bootstrap)
			if err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to migrate cluster: %w", err)
			}
			res, err = r.checkCluster(ctx, bootstrap)
			if err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to check cluster: %w", err)
			}
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Running, "Cluster created successfully", true)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			if res != nil {
				return *res, nil
			}
		}
	}
	if bootstrap.Spec.SkipClusterCreation && bootstrap.Spec.MoveCluster {
		if !bootstrap.Status.CapiClusterStatus.Ready {
			res, err := r.checkCluster(ctx, bootstrap)
			if err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to check CAPI cluster: %w", err)
			}
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Creating, "Creating cluster", false)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			if res != nil {
				return *res, nil
			}
		}

		if !bootstrap.Status.Ready {
			if err := r.moveNamespace(ctx, bootstrap); err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to move CAPI objects: %w", err)
			}
			if err := r.updateStatus(ctx, bootstrap, bv1alpha1.Running, "Cluster created successfully", true); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	if !bootstrap.Spec.SkipClusterCreation {
		if !bootstrap.Status.CapiClusterStatus.Ready {
			res, err := r.reconcileCluster(ctx, bootstrap)
			if err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to reconcile CAPI cluster: %w", err)
			}
			updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Creating, "Creating cluster", false)
			if updateErr != nil {
				return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
			}
			if res != nil {
				return *res, nil
			}
		}

		if !bootstrap.Status.Ready {
			if err := r.moveNamespace(ctx, bootstrap); err != nil {
				updateErr := r.updateStatus(ctx, bootstrap, bv1alpha1.Error, err.Error(), false)
				if updateErr != nil {
					return ctrl.Result{}, fmt.Errorf("failed to set the bootstrap error: %w", updateErr)
				}
				return ctrl.Result{}, fmt.Errorf("failed to move CAPI objects: %w", err)
			}
			if err := r.updateStatus(ctx, bootstrap, bv1alpha1.Running, "Cluster created successfully", true); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) updateStatus(ctx context.Context, bootstrap *bv1alpha1.Bootstrap, phase bv1alpha1.ComponentPhase, message string, ready bool) error {
	err := helper.UpdateBootstrapStatus(ctx, r.Client, bootstrap, func(c *bv1alpha1.Bootstrap) {
		c.Status.Message = message
		c.Status.Phase = phase
		c.Status.Ready = ready

	})
	if err != nil {
		return fmt.Errorf("failed to set error status on bootstrap to: errorMessage=%q. Could not update cluster: %w", message, err)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bv1alpha1.Bootstrap{}).
		Owns(&appsv1.Deployment{}).
		Watches(
			&source.Kind{Type: &clusterapiv1beta1.Cluster{}},
			handler.EnqueueRequestsFromMapFunc(r.findClusterObject),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		// Watches(
		// 	&source.Kind{Type: &corev1.ConfigMap{}},
		// 	handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
		// 	builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		// ).
		Complete(r)
}

func (r *Reconciler) findClusterObject(cluster client.Object) []reconcile.Request {
	attachedBootstraps := &bv1alpha1.BootstrapList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(clusterNameField, cluster.GetName()),
		Namespace:     cluster.GetNamespace(),
	}
	err := r.List(context.TODO(), attachedBootstraps, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(attachedBootstraps.Items))
	for i, item := range attachedBootstraps.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}
	return requests
}
