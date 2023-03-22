// This file is generated. DO NOT EDIT.
package reconciling

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
)

// NamespaceCreator defines an interface to create/update Namespaces
type NamespaceCreator = func(existing *corev1.Namespace) (*corev1.Namespace, error)

// NamedNamespaceCreatorGetter returns the name of the resource and the corresponding creator function
type NamedNamespaceCreatorGetter = func() (name string, create NamespaceCreator)

// NamespaceObjectWrapper adds a wrapper so the NamespaceCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func NamespaceObjectWrapper(create NamespaceCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*corev1.Namespace))
		}
		return create(&corev1.Namespace{})
	}
}

// ReconcileNamespaces will create and update the Namespaces coming from the passed NamespaceCreator slice
func ReconcileNamespaces(ctx context.Context, namedGetters []NamedNamespaceCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := NamespaceObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &corev1.Namespace{}, false); err != nil {
			return fmt.Errorf("failed to ensure Namespace %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ServiceCreator defines an interface to create/update Services
type ServiceCreator = func(existing *corev1.Service) (*corev1.Service, error)

// NamedServiceCreatorGetter returns the name of the resource and the corresponding creator function
type NamedServiceCreatorGetter = func() (name string, create ServiceCreator)

// ServiceObjectWrapper adds a wrapper so the ServiceCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ServiceObjectWrapper(create ServiceCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*corev1.Service))
		}
		return create(&corev1.Service{})
	}
}

// ReconcileServices will create and update the Services coming from the passed ServiceCreator slice
func ReconcileServices(ctx context.Context, namedGetters []NamedServiceCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ServiceObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &corev1.Service{}, false); err != nil {
			return fmt.Errorf("failed to ensure Service %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// SecretCreator defines an interface to create/update Secrets
type SecretCreator = func(existing *corev1.Secret) (*corev1.Secret, error)

// NamedSecretCreatorGetter returns the name of the resource and the corresponding creator function
type NamedSecretCreatorGetter = func() (name string, create SecretCreator)

// SecretObjectWrapper adds a wrapper so the SecretCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func SecretObjectWrapper(create SecretCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*corev1.Secret))
		}
		return create(&corev1.Secret{})
	}
}

// ReconcileSecrets will create and update the Secrets coming from the passed SecretCreator slice
func ReconcileSecrets(ctx context.Context, namedGetters []NamedSecretCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := SecretObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &corev1.Secret{}, false); err != nil {
			return fmt.Errorf("failed to ensure Secret %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ConfigMapCreator defines an interface to create/update ConfigMaps
type ConfigMapCreator = func(existing *corev1.ConfigMap) (*corev1.ConfigMap, error)

// NamedConfigMapCreatorGetter returns the name of the resource and the corresponding creator function
type NamedConfigMapCreatorGetter = func() (name string, create ConfigMapCreator)

// ConfigMapObjectWrapper adds a wrapper so the ConfigMapCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ConfigMapObjectWrapper(create ConfigMapCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*corev1.ConfigMap))
		}
		return create(&corev1.ConfigMap{})
	}
}

// ReconcileConfigMaps will create and update the ConfigMaps coming from the passed ConfigMapCreator slice
func ReconcileConfigMaps(ctx context.Context, namedGetters []NamedConfigMapCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ConfigMapObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &corev1.ConfigMap{}, false); err != nil {
			return fmt.Errorf("failed to ensure ConfigMap %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ServiceAccountCreator defines an interface to create/update ServiceAccounts
type ServiceAccountCreator = func(existing *corev1.ServiceAccount) (*corev1.ServiceAccount, error)

// NamedServiceAccountCreatorGetter returns the name of the resource and the corresponding creator function
type NamedServiceAccountCreatorGetter = func() (name string, create ServiceAccountCreator)

// ServiceAccountObjectWrapper adds a wrapper so the ServiceAccountCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ServiceAccountObjectWrapper(create ServiceAccountCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*corev1.ServiceAccount))
		}
		return create(&corev1.ServiceAccount{})
	}
}

// ReconcileServiceAccounts will create and update the ServiceAccounts coming from the passed ServiceAccountCreator slice
func ReconcileServiceAccounts(ctx context.Context, namedGetters []NamedServiceAccountCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ServiceAccountObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &corev1.ServiceAccount{}, false); err != nil {
			return fmt.Errorf("failed to ensure ServiceAccount %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// DeploymentCreator defines an interface to create/update Deployments
type DeploymentCreator = func(existing *appsv1.Deployment) (*appsv1.Deployment, error)

// NamedDeploymentCreatorGetter returns the name of the resource and the corresponding creator function
type NamedDeploymentCreatorGetter = func() (name string, create DeploymentCreator)

// DeploymentObjectWrapper adds a wrapper so the DeploymentCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func DeploymentObjectWrapper(create DeploymentCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*appsv1.Deployment))
		}
		return create(&appsv1.Deployment{})
	}
}

// ReconcileDeployments will create and update the Deployments coming from the passed DeploymentCreator slice
func ReconcileDeployments(ctx context.Context, namedGetters []NamedDeploymentCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		create = DefaultDeployment(create)
		createObject := DeploymentObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &appsv1.Deployment{}, false); err != nil {
			return fmt.Errorf("failed to ensure Deployment %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ClusterRoleBindingCreator defines an interface to create/update ClusterRoleBindings
type ClusterRoleBindingCreator = func(existing *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)

// NamedClusterRoleBindingCreatorGetter returns the name of the resource and the corresponding creator function
type NamedClusterRoleBindingCreatorGetter = func() (name string, create ClusterRoleBindingCreator)

// ClusterRoleBindingObjectWrapper adds a wrapper so the ClusterRoleBindingCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ClusterRoleBindingObjectWrapper(create ClusterRoleBindingCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*rbacv1.ClusterRoleBinding))
		}
		return create(&rbacv1.ClusterRoleBinding{})
	}
}

// ReconcileClusterRoleBindings will create and update the ClusterRoleBindings coming from the passed ClusterRoleBindingCreator slice
func ReconcileClusterRoleBindings(ctx context.Context, namedGetters []NamedClusterRoleBindingCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ClusterRoleBindingObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &rbacv1.ClusterRoleBinding{}, false); err != nil {
			return fmt.Errorf("failed to ensure ClusterRoleBinding %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ClusterRoleCreator defines an interface to create/update ClusterRoles
type ClusterRoleCreator = func(existing *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)

// NamedClusterRoleCreatorGetter returns the name of the resource and the corresponding creator function
type NamedClusterRoleCreatorGetter = func() (name string, create ClusterRoleCreator)

// ClusterRoleObjectWrapper adds a wrapper so the ClusterRoleCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ClusterRoleObjectWrapper(create ClusterRoleCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*rbacv1.ClusterRole))
		}
		return create(&rbacv1.ClusterRole{})
	}
}

// ReconcileClusterRoles will create and update the ClusterRoles coming from the passed ClusterRoleCreator slice
func ReconcileClusterRoles(ctx context.Context, namedGetters []NamedClusterRoleCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ClusterRoleObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &rbacv1.ClusterRole{}, false); err != nil {
			return fmt.Errorf("failed to ensure ClusterRole %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// RoleCreator defines an interface to create/update Roles
type RoleCreator = func(existing *rbacv1.Role) (*rbacv1.Role, error)

// NamedRoleCreatorGetter returns the name of the resource and the corresponding creator function
type NamedRoleCreatorGetter = func() (name string, create RoleCreator)

// RoleObjectWrapper adds a wrapper so the RoleCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func RoleObjectWrapper(create RoleCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*rbacv1.Role))
		}
		return create(&rbacv1.Role{})
	}
}

// ReconcileRoles will create and update the Roles coming from the passed RoleCreator slice
func ReconcileRoles(ctx context.Context, namedGetters []NamedRoleCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := RoleObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &rbacv1.Role{}, false); err != nil {
			return fmt.Errorf("failed to ensure Role %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// RoleBindingCreator defines an interface to create/update RoleBindings
type RoleBindingCreator = func(existing *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)

// NamedRoleBindingCreatorGetter returns the name of the resource and the corresponding creator function
type NamedRoleBindingCreatorGetter = func() (name string, create RoleBindingCreator)

// RoleBindingObjectWrapper adds a wrapper so the RoleBindingCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func RoleBindingObjectWrapper(create RoleBindingCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*rbacv1.RoleBinding))
		}
		return create(&rbacv1.RoleBinding{})
	}
}

// ReconcileRoleBindings will create and update the RoleBindings coming from the passed RoleBindingCreator slice
func ReconcileRoleBindings(ctx context.Context, namedGetters []NamedRoleBindingCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := RoleBindingObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &rbacv1.RoleBinding{}, false); err != nil {
			return fmt.Errorf("failed to ensure RoleBinding %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// CertificateCreator defines an interface to create/update Certificates
type CertificateCreator = func(existing *certmanagerv1.Certificate) (*certmanagerv1.Certificate, error)

// NamedCertificateCreatorGetter returns the name of the resource and the corresponding creator function
type NamedCertificateCreatorGetter = func() (name string, create CertificateCreator)

// CertificateObjectWrapper adds a wrapper so the CertificateCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func CertificateObjectWrapper(create CertificateCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*certmanagerv1.Certificate))
		}
		return create(&certmanagerv1.Certificate{})
	}
}

// ReconcileCertificates will create and update the Certificates coming from the passed CertificateCreator slice
func ReconcileCertificates(ctx context.Context, namedGetters []NamedCertificateCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := CertificateObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &certmanagerv1.Certificate{}, false); err != nil {
			return fmt.Errorf("failed to ensure Certificate %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// IssuerCreator defines an interface to create/update Issuers
type IssuerCreator = func(existing *certmanagerv1.Issuer) (*certmanagerv1.Issuer, error)

// NamedIssuerCreatorGetter returns the name of the resource and the corresponding creator function
type NamedIssuerCreatorGetter = func() (name string, create IssuerCreator)

// IssuerObjectWrapper adds a wrapper so the IssuerCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func IssuerObjectWrapper(create IssuerCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*certmanagerv1.Issuer))
		}
		return create(&certmanagerv1.Issuer{})
	}
}

// ReconcileIssuers will create and update the Issuers coming from the passed IssuerCreator slice
func ReconcileIssuers(ctx context.Context, namedGetters []NamedIssuerCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := IssuerObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &certmanagerv1.Issuer{}, false); err != nil {
			return fmt.Errorf("failed to ensure Issuer %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ValidatingWebhookConfigurationCreator defines an interface to create/update ValidatingWebhookConfigurations
type ValidatingWebhookConfigurationCreator = func(existing *admissionregistrationv1.ValidatingWebhookConfiguration) (*admissionregistrationv1.ValidatingWebhookConfiguration, error)

// NamedValidatingWebhookConfigurationCreatorGetter returns the name of the resource and the corresponding creator function
type NamedValidatingWebhookConfigurationCreatorGetter = func() (name string, create ValidatingWebhookConfigurationCreator)

// ValidatingWebhookConfigurationObjectWrapper adds a wrapper so the ValidatingWebhookConfigurationCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ValidatingWebhookConfigurationObjectWrapper(create ValidatingWebhookConfigurationCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*admissionregistrationv1.ValidatingWebhookConfiguration))
		}
		return create(&admissionregistrationv1.ValidatingWebhookConfiguration{})
	}
}

// ReconcileValidatingWebhookConfigurations will create and update the ValidatingWebhookConfigurations coming from the passed ValidatingWebhookConfigurationCreator slice
func ReconcileValidatingWebhookConfigurations(ctx context.Context, namedGetters []NamedValidatingWebhookConfigurationCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ValidatingWebhookConfigurationObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &admissionregistrationv1.ValidatingWebhookConfiguration{}, false); err != nil {
			return fmt.Errorf("failed to ensure ValidatingWebhookConfiguration %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// BootstrapProviderCreator defines an interface to create/update BootstrapProviders
type BootstrapProviderCreator = func(existing *clusterapioperator.BootstrapProvider) (*clusterapioperator.BootstrapProvider, error)

// NamedBootstrapProviderCreatorGetter returns the name of the resource and the corresponding creator function
type NamedBootstrapProviderCreatorGetter = func() (name string, create BootstrapProviderCreator)

// BootstrapProviderObjectWrapper adds a wrapper so the BootstrapProviderCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func BootstrapProviderObjectWrapper(create BootstrapProviderCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*clusterapioperator.BootstrapProvider))
		}
		return create(&clusterapioperator.BootstrapProvider{})
	}
}

// ReconcileBootstrapProviders will create and update the BootstrapProviders coming from the passed BootstrapProviderCreator slice
func ReconcileBootstrapProviders(ctx context.Context, namedGetters []NamedBootstrapProviderCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := BootstrapProviderObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &clusterapioperator.BootstrapProvider{}, false); err != nil {
			return fmt.Errorf("failed to ensure BootstrapProvider %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ControlPlaneProviderCreator defines an interface to create/update ControlPlaneProviders
type ControlPlaneProviderCreator = func(existing *clusterapioperator.ControlPlaneProvider) (*clusterapioperator.ControlPlaneProvider, error)

// NamedControlPlaneProviderCreatorGetter returns the name of the resource and the corresponding creator function
type NamedControlPlaneProviderCreatorGetter = func() (name string, create ControlPlaneProviderCreator)

// ControlPlaneProviderObjectWrapper adds a wrapper so the ControlPlaneProviderCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func ControlPlaneProviderObjectWrapper(create ControlPlaneProviderCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*clusterapioperator.ControlPlaneProvider))
		}
		return create(&clusterapioperator.ControlPlaneProvider{})
	}
}

// ReconcileControlPlaneProviders will create and update the ControlPlaneProviders coming from the passed ControlPlaneProviderCreator slice
func ReconcileControlPlaneProviders(ctx context.Context, namedGetters []NamedControlPlaneProviderCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := ControlPlaneProviderObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &clusterapioperator.ControlPlaneProvider{}, false); err != nil {
			return fmt.Errorf("failed to ensure ControlPlaneProvider %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// InfrastructureProviderCreator defines an interface to create/update InfrastructureProviders
type InfrastructureProviderCreator = func(existing *clusterapioperator.InfrastructureProvider) (*clusterapioperator.InfrastructureProvider, error)

// NamedInfrastructureProviderCreatorGetter returns the name of the resource and the corresponding creator function
type NamedInfrastructureProviderCreatorGetter = func() (name string, create InfrastructureProviderCreator)

// InfrastructureProviderObjectWrapper adds a wrapper so the InfrastructureProviderCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func InfrastructureProviderObjectWrapper(create InfrastructureProviderCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*clusterapioperator.InfrastructureProvider))
		}
		return create(&clusterapioperator.InfrastructureProvider{})
	}
}

// ReconcileInfrastructureProviders will create and update the InfrastructureProviders coming from the passed InfrastructureProviderCreator slice
func ReconcileInfrastructureProviders(ctx context.Context, namedGetters []NamedInfrastructureProviderCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := InfrastructureProviderObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &clusterapioperator.InfrastructureProvider{}, false); err != nil {
			return fmt.Errorf("failed to ensure InfrastructureProvider %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// CoreProviderCreator defines an interface to create/update CoreProviders
type CoreProviderCreator = func(existing *clusterapioperator.CoreProvider) (*clusterapioperator.CoreProvider, error)

// NamedCoreProviderCreatorGetter returns the name of the resource and the corresponding creator function
type NamedCoreProviderCreatorGetter = func() (name string, create CoreProviderCreator)

// CoreProviderObjectWrapper adds a wrapper so the CoreProviderCreator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func CoreProviderObjectWrapper(create CoreProviderCreator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*clusterapioperator.CoreProvider))
		}
		return create(&clusterapioperator.CoreProvider{})
	}
}

// ReconcileCoreProviders will create and update the CoreProviders coming from the passed CoreProviderCreator slice
func ReconcileCoreProviders(ctx context.Context, namedGetters []NamedCoreProviderCreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
		createObject := CoreProviderObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &clusterapioperator.CoreProvider{}, false); err != nil {
			return fmt.Errorf("failed to ensure CoreProvider %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}
