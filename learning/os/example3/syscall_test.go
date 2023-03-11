package example3_test

import (
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestSocketpair(t *testing.T) {
	t.Parallel()

	testData := []byte("Hello world!")
	fds, err := syscall.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	require.NoError(t, err)

	file1 := os.NewFile(uintptr(fds[0]), "file1")

	t.Cleanup(func() {
		err = file1.Close()
		require.NoError(t, err)
	})

	file2 := os.NewFile(uintptr(fds[1]), "file2")

	t.Cleanup(func() {
		err = file2.Close()
		require.NoError(t, err)
	})

	conn1, err := net.FileConn(file1)

	t.Cleanup(func() {
		err = conn1.Close()
		require.NoError(t, err)
	})

	conn2, err := net.FileConn(file2)

	t.Cleanup(func() {
		err = conn2.Close()
		require.NoError(t, err)
	})

	_, err = conn1.Write(testData)
	require.NoError(t, err)

	buf := make([]byte, len(testData))
	_, err = conn2.Read(buf)

	assert.Equal(t, testData, buf)
}
