package example1_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/k8s/example1"
	clientSet "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/clientset/versioned"
	informers "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/informers/externalversions"
	"github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/util"
	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd"
	"k8s.io/kubectl/pkg/cmd/plugin"
)

const (
	fooName               = "test.foo"
	deploymentName        = "test.deployment"
	crdConfigPath         = "./crd-status-subresource.yaml"
	namespace             = "default"
	kubeconfigPath        = "../../../docker/k8s/kubeconfig-dev"
	masterURLEnvVar       = "KUBE_API_SERVER_URL"
	reSyncDuration        = time.Second * 30
	replicas        int32 = 1
	workersCount          = 2
	delaySeconds          = 10
)

func TestController_Run(t *testing.T) {
	t.Parallel()

	outputBuf := new(utils.ThreadSafeBuffer)

	mu := sync.Mutex{}
	mu.Lock()
	logSetUp(t, outputBuf)

	masterURL, ok := os.LookupEnv(masterURLEnvVar)
	require.True(t, ok)

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	require.NoError(t, err)

	defaultConfigFlags := genericclioptions.
		NewConfigFlags(true).
		WithDeprecatedPasswordFlag().
		WithDiscoveryBurst(300).
		WithDiscoveryQPS(50.0)
	defaultConfigFlags.KubeConfig = kubeconfigPtr()

	kubeClient, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	exampleClient, err := clientSet.NewForConfig(cfg)
	require.NoError(t, err)

	kubeInformerFactory := kubeInformers.NewSharedInformerFactory(kubeClient, reSyncDuration)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, reSyncDuration)

	controller := example1.NewController(
		kubeClient,
		exampleClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		exampleInformerFactory.Samplecontroller().V1alpha1().Foos(),
	)

	err = addCRD(defaultConfigFlags)
	require.NoError(t, err)
	time.Sleep(time.Second * delaySeconds)

	_, err = exampleClient.SamplecontrollerV1alpha1().Foos(namespace).Create(
		context.Background(),
		util.NewFoo(replicas, namespace, fooName, deploymentName),
		util.NewCreateOptions(),
	)
	require.NoError(t, err)
	time.Sleep(time.Second * delaySeconds)

	defer cleanUp(t, kubeClient, exampleClient, defaultConfigFlags)

	stopCh := make(chan struct{})
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)

	go func() {
		err = controller.Run(workersCount, stopCh)
		assert.NoError(t, err)
	}()

	time.Sleep(time.Second * delaySeconds)

	assert.True(t, strings.Contains(outputBuf.String(), "Successfully synced 'default/test.foo'"))

	stopCh <- struct{}{}

	klog.Flush()
}

func logSetUp(t *testing.T, outputBuf io.Writer) {
	t.Helper()

	klog.InitFlags(nil)

	err := flag.Set("logtostderr", "false")
	require.NoError(t, err)

	err = flag.Set("alsologtostderr", "false")
	require.NoError(t, err)

	flag.Parse()
	klog.SetOutput(outputBuf)
}

func kubeconfigPtr() *string {
	kubeconfigTmp := kubeconfigPath

	return &kubeconfigTmp
}

func addCRD(defaultConfigFlags *genericclioptions.ConfigFlags) error {
	command := cmd.NewDefaultKubectlCommandWithArgs(cmd.KubectlOptions{
		PluginHandler: cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     nil,
		ConfigFlags:   defaultConfigFlags,
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
	command.SetArgs([]string{"apply", "-f", crdConfigPath})

	if err := command.Execute(); err != nil {
		return fmt.Errorf("command execution error: %w", err)
	}

	return nil
}

func deleteCRD(defaultConfigFlags *genericclioptions.ConfigFlags) error {
	command := cmd.NewDefaultKubectlCommandWithArgs(cmd.KubectlOptions{
		PluginHandler: cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     nil,
		ConfigFlags:   defaultConfigFlags,
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
	command.SetArgs([]string{"delete", "-f", crdConfigPath})

	if err := command.Execute(); err != nil {
		return fmt.Errorf("command execution error: %w", err)
	}

	return nil
}

func cleanUp(
	t *testing.T,
	kubeClient kubernetes.Interface,
	exampleClient clientSet.Interface,
	defaultConfigFlags *genericclioptions.ConfigFlags,
) {
	t.Helper()

	err := kubeClient.AppsV1().Deployments(namespace).Delete(
		context.Background(),
		deploymentName,
		util.NewDeletionOptions(),
	)
	require.NoError(t, err)

	err = exampleClient.SamplecontrollerV1alpha1().Foos(namespace).Delete(
		context.Background(),
		fooName,
		util.NewDeletionOptions(),
	)
	require.NoError(t, err)

	err = deleteCRD(defaultConfigFlags)
	require.NoError(t, err)
}
