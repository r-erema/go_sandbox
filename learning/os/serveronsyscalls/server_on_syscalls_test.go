package serveronsyscalls_test

import (
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/os/serveronsyscalls"
	"github.com/reiver/go-telnet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Parallel()

	serverIncomingData := make(chan []byte)

	go func() {
		err := serveronsyscalls.Server("127.0.0.1", 7777, serverIncomingData)
		assert.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 100)

	conn, err := telnet.DialTo("127.0.0.1:7777")
	require.NoError(t, err)

	data := []byte("TEST DATA")
	_, err = conn.Write(data)
	require.NoError(t, err)

	assert.Equal(t, data, <-serverIncomingData)
}
