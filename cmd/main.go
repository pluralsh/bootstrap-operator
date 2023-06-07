package main

import (
	"flag"
	"os"

	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clusterapioperator "sigs.k8s.io/cluster-api-operator/api/v1alpha1"
	awsinfrastructure "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	awscontrolplane "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	awsmachinepool "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	azurecontroleplane "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	gcpclusterapi "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterapi "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"
	"github.com/pluralsh/bootstrap-operator/pkg/controller"
	"github.com/pluralsh/bootstrap-operator/pkg/log"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = log.Logger
)

func init() {
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	utilruntime.Must(clusterapiexp.AddToScheme(scheme))
	utilruntime.Must(awsmachinepool.AddToScheme(scheme))
	utilruntime.Must(awsinfrastructure.AddToScheme(scheme))
	utilruntime.Must(awscontrolplane.AddToScheme(scheme))
	utilruntime.Must(azurecontroleplane.AddToScheme(scheme))
	utilruntime.Must(clusterapi.AddToScheme(scheme))
	utilruntime.Must(clusterapioperator.AddToScheme(scheme))
	utilruntime.Must(admissionregistrationv1.AddToScheme(scheme))
	utilruntime.Must(certv1.AddToScheme(scheme))
	utilruntime.Must(rbacv1.AddToScheme(scheme))
	utilruntime.Must(bv1alpha1.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(gcpclusterapi.AddToScheme(scheme))
	utilruntime.Must(cmapi.AddToScheme(scheme))

}

func main() {
	var enableLeaderElection bool
	var metricsAddr string
	var probeAddr string
	var namespace string
	var kubeconfig string
	flag.StringVar(&namespace, "namespace", "default", "The namespace operator runs in")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	config := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                 scheme,
		Namespace:              namespace,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "1237ab00.plural.sh",
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: probeAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		setupLog.Error(err, "unable to create kubernetes clientset")
		os.Exit(1)
	}

	if err = (&controller.Reconciler{
		Client:     mgr.GetClient(),
		KubeClient: clientset,
		Namespace:  namespace,
		Scheme:     scheme,
		Kubeconfig: kubeconfig,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "bootstrap")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()
	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

}
