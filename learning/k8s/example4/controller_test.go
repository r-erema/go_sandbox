package example4_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/k8s/example4/pkg"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	ingressManifestPath    = "ingress.yaml"
	certPath               = "../../../docker/k8s/common_cert_for_all.crt"
	certKey                = "../../../docker/k8s/common_cert_key_for_all.key"
	backendServiceResponse = "Backend service has been reached"
)

func TestController(t *testing.T) { //nolint: paralleltest
	defaultConfigFlags := test.CLIConfigFlags(t)

	err := test.RunKubectlCommand(
		defaultConfigFlags,
		[]string{"create", "secret", "tls", "ingress-tls", "--key", certKey, "--cert", certPath},
	)
	require.NoError(t, err)

	defer func() {
		err = test.RunKubectlCommand(defaultConfigFlags, []string{"delete", "secret", "ingress-tls"})
		require.NoError(t, err)
	}()

	err = test.RunKubectlCommand(defaultConfigFlags, []string{"apply", "-f", ingressManifestPath})
	require.NoError(t, err)

	defer func() {
		err = test.RunKubectlCommand(defaultConfigFlags, []string{"delete", "-f", ingressManifestPath})
		require.NoError(t, err)
	}()

	cfg, err := clientcmd.BuildConfigFromFlags("", *test.KubeConfigPtr(t))
	require.NoError(t, err)

	kubeClientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	server := pkg.NewServer(pkg.WithTLSPort(4040))

	watcher := pkg.NewWatcher(kubeClientset, func(payload *pkg.Payload) {
		server.Update(payload)
	})

	go func() {
		err = server.Run(context.Background())
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)

	go func() {
		err = watcher.Run(context.Background())
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)

	go func() {
		srv := http.Server{
			Addr: fmt.Sprintf("%s:%d", "0.0.0.0", 5050),
			Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, err = writer.Write([]byte(backendServiceResponse))
				require.NoError(t, err)
			}),
			ReadHeaderTimeout: time.Second,
		}

		err = srv.ListenAndServeTLS(certPath, certKey)
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://localhost:4040", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, []byte(backendServiceResponse), bodyBytes)
}
