package example1_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

const (
	runParentProcessArg    = "run-process-workload"
	childFileDescriptorEnv = "CHILD_FD"
)

func TestOpenExtraFilesCmdFiles(t *testing.T) {
	t.Parallel()

	file, err := os.Open("./test_data/test_text_data")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = file.Close()
		require.NoError(t, err)
	})

	fileDescriptor, readBytesOffset := "3", "5"

	cmd := exec.Command("./read_by_file_descriptor", fileDescriptor, readBytesOffset)
	cmd.ExtraFiles = []*os.File{file}

	output, err := cmd.Output()
	require.NoError(t, err)

	assert.Equal(t, []byte("12345"), bytes.Trim(output, "\n"))
}

func TestUsingDescriptorInOtherProcess(t *testing.T) {
	t.Parallel()

	if os.Args[1] == runParentProcessArg {
		processWorkload(t)

		return
	}

	parentSock, childSock := sockPair(t)
	t.Cleanup(func() {
		err := parentSock.Close()
		require.NoError(t, err)
		err = childSock.Close()
		require.NoError(t, err)
	})

	cmd := exec.Command("/proc/self/exe", runParentProcessArg)
	cmd.ExtraFiles = append(cmd.ExtraFiles, childSock)
	cmd.Env = []string{
		fmt.Sprintf("%s=%d", childFileDescriptorEnv, 3),
	}

	err := cmd.Start()
	require.NoError(t, err)

	err = <-waiter(t, parentSock)
	assert.NoError(t, err)
}

func waiter(t *testing.T, reader io.Reader) chan error {
	t.Helper()

	channel := make(chan error, 1)

	go func() {
		buf := make([]byte, 1)
		_, err := reader.Read(buf)
		require.NoError(t, err)
		require.Equal(t, []byte("1"), buf)
		channel <- nil
	}()

	return channel
}

func processWorkload(t *testing.T) {
	t.Helper()

	childFDEnv := os.Getenv(childFileDescriptorEnv)
	fileDescriptor, err := strconv.Atoi(childFDEnv)
	require.NoError(t, err)

	file := os.NewFile(uintptr(fileDescriptor), "")

	t.Cleanup(func() {
		err = file.Close()
		require.NoError(t, err)
	})

	_, err = file.WriteString("1")
	require.NoError(t, err)
}

func sockPair(t *testing.T) (sock1, sock2 *os.File) { //nolint:nonamedreturns
	t.Helper()

	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	require.NoError(t, err)

	return os.NewFile(uintptr(fds[1]), ""), os.NewFile(uintptr(fds[0]), "")
}
