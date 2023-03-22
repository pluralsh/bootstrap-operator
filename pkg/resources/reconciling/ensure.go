package reconciling

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/pluralsh/bootstrap-operator/pkg/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run ../../../codegen/reconcile/main.go

// ObjectCreator defines an interface to create/update a ctrlruntimeclient.Object.
type ObjectCreator = func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error)

// ObjectModifier is a wrapper function which modifies the object which gets returned by the passed in ObjectCreator.
type ObjectModifier func(create ObjectCreator) ObjectCreator

func createWithNamespace(rawcreate ObjectCreator, namespace string) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		obj, err := rawcreate(existing)
		if err != nil {
			return nil, err
		}
		obj.(metav1.Object).SetNamespace(namespace)
		return obj, nil
	}
}

func createWithName(rawcreate ObjectCreator, name string) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		obj, err := rawcreate(existing)
		if err != nil {
			return nil, err
		}
		obj.(metav1.Object).SetName(name)
		return obj, nil
	}
}

func objectLogger(obj ctrlruntimeclient.Object) *zap.SugaredLogger {
	// make sure we handle objects with broken typeMeta and still create a nice-looking kind name
	logger := log.Logger.With("kind", reflect.TypeOf(obj).Elem())
	if ns := obj.GetNamespace(); ns != "" {
		logger = logger.With("namespace", ns)
	}

	// ensure name comes after namespace
	return logger.With("name", obj.GetName())
}

// EnsureNamedObject will generate the Object with the passed create function & create or update it in Kubernetes if necessary.
func EnsureNamedObject(ctx context.Context, namespacedName types.NamespacedName, rawcreate ObjectCreator, client ctrlruntimeclient.Client, emptyObject ctrlruntimeclient.Object, requiresRecreate bool) error {
	// A wrapper to ensure we always set the Namespace and Name. This is useful as we call create twice
	create := createWithNamespace(rawcreate, namespacedName.Namespace)
	create = createWithName(create, namespacedName.Name)

	exists := true
	existingObject := emptyObject.DeepCopyObject().(ctrlruntimeclient.Object)
	if err := client.Get(ctx, namespacedName, existingObject); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get Object(%T): %w", existingObject, err)
		}
		exists = false
	}

	if exists {
		return nil
	}

	// Object does not exist in lister -> Create the Object
	obj, err := create(emptyObject)
	if err != nil {
		return fmt.Errorf("failed to generate object: %w", err)
	}
	if err := client.Create(ctx, obj); err != nil {
		return fmt.Errorf("failed to create %T '%s': %w", obj, namespacedName.String(), err)
	}
	// Wait until the object exists in the cache
	createdObjectIsInCache := WaitUntilObjectExistsInCacheConditionFunc(ctx, client, objectLogger(obj), namespacedName, obj)
	err = wait.PollImmediate(10*time.Millisecond, 10*time.Second, createdObjectIsInCache)
	if err != nil {
		return fmt.Errorf("failed waiting for the cache to contain our newly created object: %w", err)
	}

	objectLogger(obj).Info("created new resource")
	return nil

}

func WaitUntilObjectExistsInCacheConditionFunc(
	ctx context.Context,
	client ctrlruntimeclient.Client,
	log *zap.SugaredLogger,
	namespacedName types.NamespacedName,
	obj ctrlruntimeclient.Object,
) wait.ConditionFunc {
	return func() (bool, error) {
		newObj := obj.DeepCopyObject().(ctrlruntimeclient.Object)
		if err := client.Get(ctx, namespacedName, newObj); err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}

			log.Errorw("failed retrieving object while waiting for the cache to contain our newly created object", zap.Error(err))
			return false, nil
		}

		return true, nil
	}
}
