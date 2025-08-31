package queue_test

import (
	"syscall"
	"testing"

	"github.com/r-erema/go_sendbox/pkg/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	queueID, err := queue.Open("/test_queue", syscall.O_RDWR|syscall.O_CREAT)
	require.NoError(t, err)
	err = queue.Send(queueID, []byte("test data"), 1)
	require.NoError(t, err)
	require.True(t, queue.Exists("/test_queue"))
	require.NoError(t, err)

	data, err := queue.Receive(queueID)
	assert.Equal(t, []byte("test data"), data)
	t.Cleanup(func() {
		err = queue.Close(queueID)
		require.NoError(t, err)
	})
}
