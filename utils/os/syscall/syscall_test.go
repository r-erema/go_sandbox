package syscall_test

import (
	"os"
	"testing"

	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestReadWriteSyscall(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp(os.TempDir(), "")
	require.NoError(t, err)

	fd, err := unix.Open(tempFile.Name(), unix.O_RDWR, 0o600)
	require.NoError(t, err)

	err = syscall.WriteToFileFD(fd, []byte("hello world"))
	require.NoError(t, err)

	s, err := syscall.ReadFromFileFD(fd)
	require.NoError(t, err)
	require.Equal(t, "hello world", s)
}
