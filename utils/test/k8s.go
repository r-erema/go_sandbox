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
