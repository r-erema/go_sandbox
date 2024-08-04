package example1_test

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestGetNetworkNamespaceDescriptor(t *testing.T) {
	t.Parallel()

	ID, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)
	assert.Positive(t, ID)
}

func TestNewNetworkNamespace(t *testing.T) {
	t.Parallel()

	fdBeforeNewNS, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), syscall.Gettid())
	require.NoError(t, err)

	err = syscall.Unshare(unix.CLONE_NEWNET)
	require.NoError(t, err)

	fdAfterNewNS, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), syscall.Gettid())
	require.NoError(t, err)

	assert.NotEqual(t, fdBeforeNewNS, fdAfterNewNS)
}

func TestSetNamespaceToExtraneousProcess(t *testing.T) {
	t.Parallel()

	rootNamespaceFD, err := utils.NewNetworkNamespaceDescriptor(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)
	rootNamespaceInode, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)

	err = syscall.Unshare(unix.CLONE_NEWNET)
	require.NoError(t, err)
	newNamespacesFD, err := utils.NewNetworkNamespaceDescriptor(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)
	newNamespaceInode, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)

	assert.NotEqual(t, rootNamespaceInode, newNamespaceInode)

	returnToRootNamespace := func() {
		err = unix.Setns(int(rootNamespaceFD), unix.CLONE_NEWNET)
		require.NoError(t, err)
	}
	returnToRootNamespace()

	newNamespacesFDInNewProcess := "3"
	cmd := exec.Command("./test_data/test_bin_source/process_assigned_to_namespace/main", newNamespacesFDInNewProcess)
	cmd.ExtraFiles = append(cmd.ExtraFiles, os.NewFile(newNamespacesFD, "ns-fd"))
	output, err := cmd.Output()
	require.NoError(t, err)
	extraneousProcessNamespaceInode, err := strconv.Atoi(string(output))
	require.NoError(t, err)

	currentNamespaceInode, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), unix.Gettid())
	require.NoError(t, err)
	assert.Equal(t, currentNamespaceInode, rootNamespaceInode)
	assert.Equal(t, newNamespaceInode, extraneousProcessNamespaceInode)
}

func TestNewPIDNamespace(t *testing.T) {
	t.Parallel()

	parentPid, parentTid := os.Getpid(), syscall.Gettid()

	parentPidNS, err := utils.PIDNamespaceInodeNumber(parentPid, parentTid)
	require.NoError(t, err)

	childPid, _, errno := unix.Syscall(syscall.SYS_CLONE, uintptr(syscall.SIGCHLD|syscall.CLONE_NEWNS|syscall.CLONE_NEWPID), 0, 0)
	require.Equal(t, 0, int(errno))

	if os.Getpid() == parentPid {
		childPidNS, err := utils.PIDNamespaceInodeNumber(int(childPid), int(childPid))
		require.NoError(t, err)
		assert.NotEqual(t, uintptr(0), childPid)
		assert.NotEqual(t, parentPidNS, childPidNS)
	} else {
		assert.Equal(t, uintptr(0), childPid)
		assert.Equal(t, 1, os.Getpid())
	}
}
