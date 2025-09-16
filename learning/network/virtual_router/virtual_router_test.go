package virtual_router_test

import (
	"net"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/network/virtual_router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	t.Parallel()

	gatewayPortCh := make(chan uint16)
	r := virtual_router.NewRouter()

	go r.Run(t, gatewayPortCh)

	time.Sleep(100 * time.Millisecond)

	gatewayPort := <-gatewayPortCh

	mac1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)

	host1 := virtual_router.NewHost(t, mac1)
	host1.ConnectToGateway(gatewayPort)

	mac2, err := net.ParseMAC("AA:BB:CC:DD:EE:FF")
	require.NoError(t, err)

	host2 := virtual_router.NewHost(t, mac2)
	host2.ConnectToGateway(gatewayPort)

	host1.Send(host2.IpAddr(), []byte("hello world!"))

	rcv := host2.Receive()

	assert.Equal(t, []byte("hello world!"), rcv)
}
