package vpn

import (
	"fmt"
	"testing"

	"github.com/r-erema/go_sendbox/utils/net"
	"github.com/stretchr/testify/require"
)

func Client() (func(), error) {
	tun, err := net.SetupTun("tun5", []string{"0.0.0.0/1", "128.0.0.0/1"})
	if err != nil {
		return nil, fmt.Errorf("setup tun error: %w", err)
	}

	_ = tun

	return func() {
	}, nil
}

func TestClient(t *testing.T) {
	_, err := Client()

	t.Cleanup(func() {
		err = net.RemoveTun("tun5")
		require.NoError(t, err)
	})

	require.NoError(t, err)
}
