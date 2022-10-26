package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd"
	"k8s.io/kubectl/pkg/cmd/plugin"
)

const (
	kubeConfigPathEnvVar = "KUBECONFIG"
	discoveryBurst       = 300
	discoveryQPS         = 50.0
)

func CLIConfigFlags(t *testing.T) *genericclioptions.ConfigFlags {
	t.Helper()

	defaultConfigFlags := DefaultConfigFlags()
	defaultConfigFlags.KubeConfig = KubeConfigPtr(t)

	return defaultConfigFlags
}

func DefaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.
		NewConfigFlags(true).
		WithDeprecatedPasswordFlag().
		WithDiscoveryBurst(discoveryBurst).
		WithDiscoveryQPS(discoveryQPS)
}

func KubeConfigPtr(t *testing.T) *string {
	t.Helper()

	kubeconfigPath, ok := os.LookupEnv(kubeConfigPathEnvVar)
	require.True(t, ok)

	return &kubeconfigPath
}

func RunKubectlCommand(defaultConfigFlags *genericclioptions.ConfigFlags, args []string) error {
	command := cmd.NewDefaultKubectlCommandWithArgs(cmd.KubectlOptions{
		PluginHandler: cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     nil,
		ConfigFlags:   defaultConfigFlags,
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
	command.SetArgs(args)

	if err := command.Execute(); err != nil {
		return fmt.Errorf("command execution error: %w", err)
	}

	return nil
}

const permissions = 0o755

func PrepareKubeConfigContext(
	t *testing.T,
	kubeConfigPath,
	clusterName,
	clusterUser,
	clusterContext,
	apiServerAddr,
	certAuthority,
	certPath,
	certKeyPath string,
) {
	t.Helper()

	_, err := os.OpenFile(kubeConfigPath, os.O_WRONLY|os.O_CREATE, permissions) //nolint: nosnakecase
	require.NoError(t, err)

	kubeConfigFlag := "--kubeconfig=" + kubeConfigPath

	err = RunKubectlCommand(DefaultConfigFlags(), []string{
		"config",
		"set-credentials",
		clusterUser,
		kubeConfigFlag,
		"--client-certificate=" + certPath,
		"--client-key=" + certKeyPath,
		"--embed-certs=true",
	})
	require.NoError(t, err)

	err = RunKubectlCommand(DefaultConfigFlags(), []string{
		"config",
		"set-cluster",
		clusterName,
		kubeConfigFlag,
		"--certificate-authority=" + certAuthority,
		"--server=" + apiServerAddr,
	})
	require.NoError(t, err)

	err = RunKubectlCommand(DefaultConfigFlags(), []string{
		"config",
		"set-context",
		clusterContext,
		kubeConfigFlag,
		"--cluster=" + clusterName,
		"--user=" + clusterUser,
	})
	require.NoError(t, err)

	err = RunKubectlCommand(DefaultConfigFlags(), []string{"config", "use-context", clusterContext, kubeConfigFlag})
	require.NoError(t, err)
}
