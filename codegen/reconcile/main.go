package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func main() {
	data := struct {
		Resources []reconcileFunctionData
	}{
		Resources: []reconcileFunctionData{
			{
				ResourceName:       "Namespace",
				ImportAlias:        "corev1",
				ResourceImportPath: "k8s.io/api/core/v1",
			},
			{
				ResourceName:       "Service",
				ImportAlias:        "corev1",
				ResourceImportPath: "k8s.io/api/core/v1",
			},
			{
				ResourceName: "Secret",
				ImportAlias:  "corev1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName: "ConfigMap",
				ImportAlias:  "corev1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName: "ServiceAccount",
				ImportAlias:  "corev1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName:       "Deployment",
				ImportAlias:        "appsv1",
				ResourceImportPath: "k8s.io/api/apps/v1",
				DefaultingFunc:     "DefaultDeployment",
			},
			{
				ResourceName:       "ClusterRoleBinding",
				ImportAlias:        "rbacv1",
				ResourceImportPath: "k8s.io/api/rbac/v1",
			},
			{
				ResourceName: "ClusterRole",
				ImportAlias:  "rbacv1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName: "Role",
				ImportAlias:  "rbacv1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName: "RoleBinding",
				ImportAlias:  "rbacv1",
				// Don't specify ResourceImportPath so this block does not create a new import line in the generated code
			},
			{
				ResourceName:       "Certificate",
				ImportAlias:        "certmanagerv1",
				ResourceImportPath: "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1",
			},
			{
				ResourceName: "Issuer",
				ImportAlias:  "certmanagerv1",
			},
			{
				ResourceName:       "ValidatingWebhookConfiguration",
				ImportAlias:        "admissionregistrationv1",
				ResourceImportPath: "k8s.io/api/admissionregistration/v1",
			},
			{
				ResourceName:       "BootstrapProvider",
				ImportAlias:        "clusterapioperator",
				ResourceImportPath: "sigs.k8s.io/cluster-api-operator/api/v1alpha1",
			},
			{
				ResourceName: "ControlPlaneProvider",
				ImportAlias:  "clusterapioperator",
			},
			{
				ResourceName: "InfrastructureProvider",
				ImportAlias:  "clusterapioperator",
			},
			{
				ResourceName: "CoreProvider",
				ImportAlias:  "clusterapioperator",
			},
			{
				ResourceName:       "Cluster",
				ImportAlias:        "clusterapi",
				ResourceImportPath: "sigs.k8s.io/cluster-api/api/v1beta1",
			},
			{
				ResourceName:       "MachinePool",
				ImportAlias:        "clusterapiexp",
				ResourceImportPath: "sigs.k8s.io/cluster-api/exp/api/v1beta1",
			},
			{
				ResourceName:       "AWSManagedCluster",
				ImportAlias:        "awsinfrastructure",
				ResourceImportPath: "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2",
			},
			{
				ResourceName:       "AWSManagedMachinePool",
				ImportAlias:        "awsmachinepool",
				ResourceImportPath: "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2",
			},
			{
				ResourceName:       "AWSManagedControlPlane",
				ImportAlias:        "awscontrolplane",
				ResourceImportPath: "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2",
			},
			{
				ResourceName:       "GCPManagedCluster",
				ImportAlias:        "gcpmanagedcluster",
				ResourceImportPath: "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1",
			},
			{
				ResourceName: "GCPManagedControlPlane",
				ImportAlias:  "gcpmanagedcluster",
			},
			{
				ResourceName: "GCPManagedMachinePool",
				ImportAlias:  "gcpmanagedcluster",
			},
			{
				ResourceName:       "CustomResourceDefinition",
				ImportAlias:        "apiextensionsv1",
				ResourceImportPath: "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
			},
			{
				ResourceName:       "AzureManagedControlPlane",
				ImportAlias:        "azurecontroleplane",
				ResourceImportPath: "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1",
			},
			{
				ResourceName: "AzureManagedCluster",
				ImportAlias:  "azurecontroleplane",
			},
			{
				ResourceName: "AzureManagedMachinePool",
				ImportAlias:  "azurecontroleplane",
			},
			{
				ResourceName: "AzureClusterIdentity",
				ImportAlias:  "azurecontroleplane",
			},
		},
	}

	buf := &bytes.Buffer{}
	if err := reconcileAllTemplate.Execute(buf, data); err != nil {
		log.Fatal(err)
	}

	fmtB, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("zz_generated_reconcile.go", fmtB, 0644); err != nil {
		log.Fatal(err)
	}
}

func lowercaseFirst(str string) string {
	return strings.ToLower(string(str[0])) + str[1:]
}

var (
	reconcileAllTplFuncs = map[string]interface{}{
		"namedReconcileFunc": namedReconcileFunc,
	}
	reconcileAllTemplate = template.Must(template.New("").Funcs(reconcileAllTplFuncs).Funcs(sprig.TxtFuncMap()).Parse(`// This file is generated. DO NOT EDIT.
package reconciling

import (
	"fmt"
	"context"

	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
{{ range .Resources }}
{{- if .ResourceImportPath }}
	{{ .ImportAlias }} "{{ .ResourceImportPath }}"
{{- end }}
{{- end }}
)

{{ range .Resources }}
{{ namedReconcileFunc .ResourceName .ImportAlias .DefaultingFunc .RequiresRecreate .ResourceNamePlural .APIVersionPrefix}}
{{- end }}

`))
)

type reconcileFunctionData struct {
	ResourceName       string
	ResourceNamePlural string
	ResourceImportPath string
	ImportAlias        string
	// Optional: A defaulting func for the given object type
	// Must be defined inside the resources package
	DefaultingFunc string
	// Whether the resource must be recreated instead of updated. Required
	// e.G. for PDBs
	RequiresRecreate bool
	// Optional: adds an api version prefix to the generated functions to avoid duplication when different resources
	// have the same ResourceName
	APIVersionPrefix string
}

func namedReconcileFunc(resourceName, importAlias, defaultingFunc string, requiresRecreate bool, plural, apiVersionPrefix string) (string, error) {
	if len(plural) == 0 {
		plural = fmt.Sprintf("%ss", resourceName)
	}

	b := &bytes.Buffer{}
	err := namedReconcileFunctionTemplate.Execute(b, struct {
		ResourceName       string
		ResourceNamePlural string
		ImportAlias        string
		DefaultingFunc     string
		RequiresRecreate   bool
		APIVersionPrefix   string
	}{
		ResourceName:       resourceName,
		ResourceNamePlural: plural,
		ImportAlias:        importAlias,
		DefaultingFunc:     defaultingFunc,
		RequiresRecreate:   requiresRecreate,
		APIVersionPrefix:   apiVersionPrefix,
	})

	if err != nil {
		return "", err
	}

	return b.String(), nil
}

var (
	reconcileFunctionTplFuncs = map[string]interface{}{
		"lowercaseFirst": lowercaseFirst,
	}
)

var namedReconcileFunctionTemplate = template.Must(template.New("").Funcs(reconcileFunctionTplFuncs).Parse(`// {{ .APIVersionPrefix }}{{ .ResourceName }}Creator defines an interface to create/update {{ .ResourceName }}s
type {{ .APIVersionPrefix }}{{ .ResourceName }}Creator = func(existing *{{ .ImportAlias }}.{{ .ResourceName }}) (*{{ .ImportAlias }}.{{ .ResourceName }}, error)

// Named{{ .APIVersionPrefix }}{{ .ResourceName }}CreatorGetter returns the name of the resource and the corresponding creator function
type Named{{ .APIVersionPrefix }}{{ .ResourceName }}CreatorGetter = func() (name string, create {{ .APIVersionPrefix }}{{ .ResourceName }}Creator)

// {{ .APIVersionPrefix }}{{ .ResourceName }}ObjectWrapper adds a wrapper so the {{ .APIVersionPrefix }}{{ .ResourceName }}Creator matches ObjectCreator.
// This is needed as Go does not support function interface matching.
func {{ .APIVersionPrefix }}{{ .ResourceName }}ObjectWrapper(create {{ .APIVersionPrefix }}{{ .ResourceName }}Creator) ObjectCreator {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return create(existing.(*{{ .ImportAlias }}.{{ .ResourceName }}))
		}
		return create(&{{ .ImportAlias }}.{{ .ResourceName }}{})
	}
}

// Reconcile{{ .APIVersionPrefix }}{{ .ResourceNamePlural }} will create and update the {{ .APIVersionPrefix }}{{ .ResourceNamePlural }} coming from the passed {{ .APIVersionPrefix }}{{ .ResourceName }}Creator slice
func Reconcile{{ .APIVersionPrefix }}{{ .ResourceNamePlural }}(ctx context.Context, namedGetters []Named{{ .APIVersionPrefix }}{{ .ResourceName }}CreatorGetter, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, get := range namedGetters {
		name, create := get()
{{- if .DefaultingFunc }}
		create = {{ .DefaultingFunc }}(create)
{{- end }}
		createObject := {{ .APIVersionPrefix }}{{ .ResourceName }}ObjectWrapper(create)
		createObject = createWithNamespace(createObject, namespace)
		createObject = createWithName(createObject, name)

		for _, objectModifier := range objectModifiers {
			createObject = objectModifier(createObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, createObject, client, &{{ .ImportAlias }}.{{ .ResourceName }}{}, {{ .RequiresRecreate}}); err != nil {
			return fmt.Errorf("failed to ensure {{ .ResourceName }} %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

`))
