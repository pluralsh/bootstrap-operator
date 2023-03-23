package helper

import (
	"context"
	"reflect"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"k8s.io/client-go/util/retry"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type BootstrapPatchFunc func(bootstrap *bv1alpha1.Bootstrap)

// UpdateBootstrapStatus will attempt to patch the bootstrap status.
func UpdateBootstrapStatus(ctx context.Context, client ctrlruntimeclient.Client, bootstrap *bv1alpha1.Bootstrap, patch BootstrapPatchFunc) error {
	key := ctrlruntimeclient.ObjectKeyFromObject(bootstrap)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// fetch the current state of the cluster
		if err := client.Get(ctx, key, bootstrap); err != nil {
			return err
		}

		// modify it
		original := bootstrap.DeepCopy()
		patch(bootstrap)

		// save some work
		if reflect.DeepEqual(original.Status, bootstrap.Status) {
			return nil
		}

		// update the status
		return client.Status().Patch(ctx, bootstrap, ctrlruntimeclient.MergeFrom(original))
	})
}

// UpdateBootstrap will attempt to patch the bootstrap.
func UpdateBootstrap(ctx context.Context, client ctrlruntimeclient.Client, bootstrap *bv1alpha1.Bootstrap, patch BootstrapPatchFunc) error {
	key := ctrlruntimeclient.ObjectKeyFromObject(bootstrap)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// fetch the current state of the cluster
		if err := client.Get(ctx, key, bootstrap); err != nil {
			return err
		}

		// modify it
		original := bootstrap.DeepCopy()
		patch(bootstrap)

		// save some work
		if reflect.DeepEqual(original.Status, bootstrap.Status) {
			return nil
		}

		// update the status
		return client.Patch(ctx, bootstrap, ctrlruntimeclient.MergeFrom(original))
	})
}
