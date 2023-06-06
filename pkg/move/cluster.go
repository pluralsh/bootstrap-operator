package move

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	bv1alpha1 "github.com/pluralsh/bootstrap-operator/apis/bootstrap/v1alpha1"

	apiclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"

	coreapi "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Cluster struct {
	Namespace   string
	ClusterName string
	Bootstrap   *bv1alpha1.Bootstrap
	Scheme      *runtime.Scheme
	Client      ctrlruntimeclient.Client
	Ctx         context.Context
	Log         logr.Logger
}

func (c *Cluster) GetClient() (ctrlruntimeclient.Client, error) {
	c.Log.Info("Get bootstrap client")
	rawKubeconfig, err := c.getKubeconfig()
	if err != nil {
		return nil, err
	}
	cfg, err := clientcmd.Load(rawKubeconfig)
	if err != nil {
		c.Log.Error(err, "failed to load config")
		return nil, err
	}
	clientConfig, err := getRestConfig(cfg)
	if err != nil {
		c.Log.Error(err, "failed to get rest config")
		return nil, err
	}
	client, err := ctrlruntimeclient.New(clientConfig, ctrlruntimeclient.Options{
		Scheme: c.Scheme,
	})
	if err != nil {
		c.Log.Error(err, "failed to create client")
		return nil, err
	}
	return client, nil
}

func (c *Cluster) MoveClusterAPI() error {
	c.Log.Info("get bootstrap kubeconfig")
	k, err := c.getKubeconfig()
	if err != nil {
		return err
	}
	fileTo, err := os.CreateTemp("", "to.config")
	if err != nil {
		return err
	}

	_, err = fileTo.Write(k)
	if err != nil {
		return err
	}
	if err := fileTo.Sync(); err != nil {
		return err
	}
	c.Log.Info("bootstrap kubeconfig saved")

	var kubeconfigSecret coreapi.Secret
	if err := c.Client.Get(c.Ctx, ctrlruntimeclient.ObjectKey{Name: "kubeconfig", Namespace: c.Namespace}, &kubeconfigSecret); err != nil {
		c.Log.Error(err, "failed to get bootstrap cluster kubeconfig")
		return err
	}
	fromConfig := kubeconfigSecret.Data["value"]
	if len(fromConfig) == 0 {
		return fmt.Errorf("invalid bootstrap cluster kubeconfig, length = 0")
	}

	fileFrom, err := os.CreateTemp("", "from.config")
	if err != nil {
		return err
	}

	_, err = fileFrom.Write(fromConfig)
	if err != nil {
		return err
	}
	if err := fileFrom.Sync(); err != nil {
		return err
	}
	os.Setenv("KUBECONFIG", fileFrom.Name())

	c.Log.Info("create api client")
	client, err := apiclient.New("")
	if err != nil {
		return err
	}

	options := apiclient.MoveOptions{
		FromKubeconfig: apiclient.Kubeconfig{
			Path:    fileFrom.Name(),
			Context: "kind-bootstrap",
		},
		ToKubeconfig: apiclient.Kubeconfig{
			Path: fileTo.Name(),
		},
		Namespace: c.Namespace,
		DryRun:    false,
	}
	c.Log.Info("started moving CAPI resources ...")
	if err := client.Move(options); err != nil {
		return err
	}
	c.Log.Info("finished!")
	return nil
}

func (c *Cluster) getKubeconfig() ([]byte, error) {
	c.Log.Info("Get bootstrap kubeconfig from secret")
	var kc coreapi.Secret
	if err := c.Client.Get(c.Ctx, ctrlruntimeclient.ObjectKey{Name: fmt.Sprintf("%s-kubeconfig", c.ClusterName), Namespace: c.Namespace}, &kc); err != nil {
		return nil, err
	}

	return kc.Data["value"], nil
}

func getRestConfig(cfg *clientcmdapi.Config) (*rest.Config, error) {
	iconfig := clientcmd.NewNonInteractiveClientConfig(
		*cfg,
		"",
		&clientcmd.ConfigOverrides{},
		nil,
	)

	clientConfig, err := iconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	// Avoid blocking of the controller by increasing the QPS for user cluster interaction
	clientConfig.QPS = 120
	clientConfig.Burst = 150

	return clientConfig, nil
}
