package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	kubeConfigPathEnvVar = "KUBECONFIG"
	discoveryBurst       = 300
	discoveryQPS         = 50.0
)

func CLIConfigFlags(t *testing.T) *genericclioptions.ConfigFlags {
	t.Helper()

	defaultConfigFlags := genericclioptions.
		NewConfigFlags(true).
		WithDeprecatedPasswordFlag().
		WithDiscoveryBurst(discoveryBurst).
		WithDiscoveryQPS(discoveryQPS)
	defaultConfigFlags.KubeConfig = KubeConfigPtr(t)

	return defaultConfigFlags
}

func KubeConfigPtr(t *testing.T) *string {
	t.Helper()

	kubeconfigPath, ok := os.LookupEnv(kubeConfigPathEnvVar)
	require.True(t, ok)

	kubeconfigTmp := kubeconfigPath

	return &kubeconfigTmp
}
