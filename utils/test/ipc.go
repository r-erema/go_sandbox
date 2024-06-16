package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func SockPair(t *testing.T) (sock1, sock2 *os.File) { //nolint:nonamedreturns
	t.Helper()

	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	require.NoError(t, err)

	return os.NewFile(uintptr(fds[1]), ""), os.NewFile(uintptr(fds[0]), "")
}
