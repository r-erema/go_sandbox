package test

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func PreventSharedPropagationToRootPoint(t *testing.T, targetDir string) {
	t.Helper()

	err := syscall.Mount(targetDir, "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")
	require.NoError(t, err)
}

func MakeDirRootPoint(t *testing.T, targetDir string) {
	t.Helper()

	err := syscall.Mount(targetDir, targetDir, "", syscall.MS_BIND, "")
	require.NoError(t, err)
}

func MountFSToDirectory(t *testing.T, source, target string) {
	t.Helper()

	flags := syscall.MS_REC | syscall.MS_BIND | syscall.MS_PRIVATE
	err := syscall.Mount(source, target, "tmpfs", uintptr(flags), "mode=0755")
	require.NoError(t, err)
}
