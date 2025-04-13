package socket_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/pkg/socket"
	"github.com/r-erema/go_sendbox/pkg/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Parallel()

	clientPubKey, err := tls.Rand32Bytes()
	require.NoError(t, err)

	resp, err := socket.ClientCall("142.250.150.99", 443, clientPubKey, []byte("Hello from client"))
	require.NoError(t, err)
	assert.Equal(t, []byte("Server has got message: Hello from client"), resp)
}
