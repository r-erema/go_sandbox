package example1_test

import (
	"syscall"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkExec(t *testing.T) {
	t.Parallel()

	pid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	require.Equal(t, 0, int(errno))

	childProcess := pid == 0

	var args []string

	if childProcess {
		bin, err := syscall.BytePtrFromString("/usr/bin/echo")
		require.NoError(t, err)
		args, err := syscall.SlicePtrFromStrings([]string{"", "A child process has been run\n"})
		require.NoError(t, err)
		env, err := syscall.SlicePtrFromStrings([]string{})
		require.NoError(t, err)

		_, _, errno = syscall.Syscall(syscall.SYS_EXECVE,
			uintptr(unsafe.Pointer(bin)),
			uintptr(unsafe.Pointer(&args[0])),
			uintptr(unsafe.Pointer(&env[0])),
		)
		require.Equal(t, 0, int(errno))
	}

	assert.Empty(t, args)
}
