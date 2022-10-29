package example3_test

import (
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/k8s/example3"
	v1 "github.com/r-erema/go_sendbox/learning/k8s/example3/pkg/apis/solarsystem/v1"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	crdConfigPath = "./300-star-crd.yaml"
	crObjectPath  = "./sun-cr.yaml"
)

func TestController(t *testing.T) { //nolint: paralleltest
	defaultConfigFlags := test.CLIConfigFlags(t)

	err := test.RunKubectlCommand(defaultConfigFlags, []string{"apply", "-f", crdConfigPath})
	require.NoError(t, err)

	defer func() {
		err = test.RunKubectlCommand(defaultConfigFlags, []string{"delete", "-f", crdConfigPath})
		require.NoError(t, err)
	}()
	time.Sleep(time.Second)

	err = test.RunKubectlCommand(defaultConfigFlags, []string{"apply", "-f", crObjectPath})
	require.NoError(t, err)

	defer func() {
		err = test.RunKubectlCommand(defaultConfigFlags, []string{"delete", "-f", crObjectPath})
		require.NoError(t, err)
	}()

	syncPeriod := time.Second * 2
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		SyncPeriod: &syncPeriod,
	})
	require.NoError(t, err)

	err = v1.AddToScheme(mgr.GetScheme())
	require.NoError(t, err)

	reconciler := &example3.StarReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	star := &v1.Star{}

	err = ctrl.NewControllerManagedBy(mgr).For(star).Complete(reconciler)
	require.NoError(t, err)

	go func() {
		err = mgr.Start(ctrl.SetupSignalHandler())
		require.NoError(t, err)
	}()

	time.Sleep(time.Second * 5)

	assert.GreaterOrEqual(t, reconciler.ReconciledTimes(), uint32(2))
}
