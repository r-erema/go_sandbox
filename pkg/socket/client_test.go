package socket_test

import (
	"io"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/pkg/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Parallel()

	server := http.Server{
		Addr: ":7789",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			_, err = w.Write(slices.Concat([]byte("reply from server: "), body))
			assert.NoError(t, err)
			w.WriteHeader(http.StatusOK)
		}),
		ReadHeaderTimeout: 0,
	}

	go func() {
		err := server.ListenAndServeTLS(
			"/home/erema/h/go_sandbox/k8s/dev_environment/assets/rootCA.crt",
			"/home/erema/h/go_sandbox/k8s/dev_environment/assets/rootCA.key",
		)
		assert.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 500)

	msg := []byte(`GET / HTTP/1.1
Host: vvvc
Content-Type: application/x-www-form-urlencoded
Content-Length: 22

HTTP request test data`)

	resp, err := socket.ClientCall("localhost", 7789, msg)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "reply from server: HTTP request test data")
}
