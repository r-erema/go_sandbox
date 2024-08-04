package example2_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/containerization/example2"
	"github.com/r-erema/go_sendbox/utils"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRootFS = "./test_data/alpine-3.19.1"

	containerHostName   = "container_host_name"
	containerDomainName = "container_domain_name"
)

type (
	runCommandInContainerFn func(t *testing.T, cmd string, startInBackground bool, arguments ...string) []byte
	asserFn                 func(t *testing.T)
)

func TestContainer(t *testing.T) {
	t.Parallel()

	containerRootPath, err := filepath.Abs(testRootFS)
	require.NoError(t, err)

	prepareHelpers(t, containerRootPath)

	parentPID := os.Getpid()

	parentSocket, childSocket := test.SockPair(t)
	t.Cleanup(func() {
		err = parentSocket.Close()
		require.NoError(t, err)
	})

	parentSocketReader := bufio.NewReader(parentSocket)

	var mtx sync.Mutex

	runCommandInContainer := func(t *testing.T, cmd string, startInBackground bool, arguments ...string) []byte {
		t.Helper()

		mtx.Lock()
		defer mtx.Unlock()

		_, err = parentSocket.Write(command(t, cmd, startInBackground, arguments...))
		require.NoError(t, err)

		data, runCmdErr := parentSocketReader.ReadBytes(example2.StreamDelimiter)
		require.NoError(t, runCmdErr)

		return example2.TrimOutput(data)
	}

	containerPID, err := example2.Container(containerRootPath, containerHostName, containerDomainName, childSocket)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = example2.KillContainer(containerPID)
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		assertFn func(t *testing.T)
	}{
		{
			name:     "host and container are in different mount namespaces",
			assertFn: hostAndContainerInDifferentMountNamespacesAssertFn(parentPID, runCommandInContainer),
		},
		{
			name:     "host does not see the mount point in the container",
			assertFn: hostDoesNotSeeTheMountPointInTheContainerAssertFn(runCommandInContainer),
		},
		{
			name:     "container is able to see what was mounted in the host",
			assertFn: containerIsAbleToSeeWhatWasMountedInTheHostAssertFn(containerRootPath, runCommandInContainer),
		},
		{
			name:     "only 2 processes in the container",
			assertFn: only2ProcessesInTheContainerAssertFn(runCommandInContainer),
		},
		{
			name:     "host sees processes in the container but with different IDs",
			assertFn: hostSeesProcessesInContainerButWithDifferentIDs(runCommandInContainer),
		},
		{
			name:     "kill a process running in the container from the host",
			assertFn: killProcessRunningInContainerFromHost(runCommandInContainer),
		},
		{
			name:     "the same hostname in the host and in the container",
			assertFn: sameHostnameInHostAndInContainer(runCommandInContainer),
		},
		{
			name:     "container does not see queue of the host",
			assertFn: containerDoesNotSeeQueueOfHost(runCommandInContainer),
		},
		{
			name:     "ping the host from the container",
			assertFn: pingHostFromContainer(runCommandInContainer),
		},
		{
			name:     "limit memory by cgroup",
			assertFn: limitMemoryByCgroup(runCommandInContainer),
		},
	}

	if parentPID == os.Getpid() { // proceed with test only in parent process
		err = childSocket.Close()
		require.NoError(t, err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				tt.assertFn(t)
			})
		}
	}
}

func hostAndContainerInDifferentMountNamespacesAssertFn(parentPID int, runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		hostNamespace, err := utils.MountNamespaceInodeNumber(parentPID)
		require.NoError(t, err)
		nsLink := strings.Trim(string(runCommandInContainer(t, "readlink", false, fmt.Sprintf("/proc/%d/ns/mnt", 1))), "\n")
		containerNamespace, err := utils.ParseNSID("mnt:[", nsLink)
		require.NoError(t, err)
		assert.NotEqual(t, hostNamespace, containerNamespace)
	}
}

func hostDoesNotSeeTheMountPointInTheContainerAssertFn(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		targetPath := filepath.Join(string(os.PathSeparator), "tmp", "target_path_in_container")
		targetTestFilePath := filepath.Join(targetPath, "test_file")

		runCommandInContainer(t, "mount", false, "--verbose", "--bind", "/tmp/source_path_in_container", targetPath)
		containerOutput := string(runCommandInContainer(t, "cat", false, targetTestFilePath))
		assert.Equal(t, "test data", containerOutput)

		supposedFilePathOnHost := filepath.Join(testRootFS, targetTestFilePath)
		_, err := os.Stat(supposedFilePathOnHost)
		assert.True(t, os.IsNotExist(err))
	}
}

func containerIsAbleToSeeWhatWasMountedInTheHostAssertFn(containerRootPath string, runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		fileContent := "test data"
		source := t.TempDir()
		file, err := os.CreateTemp(source, "")
		require.NoError(t, err)
		_, err = file.WriteString(fileContent)
		require.NoError(t, err)

		target := filepath.Join(containerRootPath, "tmp", "visible_target_in_container")
		err = syscall.Mount(source, target, "", syscall.MS_BIND|syscall.MS_REC, "")
		require.NoError(t, err)
		t.Cleanup(func() {
			err = syscall.Unmount(target, syscall.MNT_FORCE|syscall.MNT_DETACH)
			require.NoError(t, err)
		})

		path := "/tmp/visible_target_in_container/" + filepath.Base(file.Name())
		containerOutput := string(runCommandInContainer(t, "cat", false, path))
		assert.Equal(t, fileContent, containerOutput)
	}
}

func only2ProcessesInTheContainerAssertFn(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		containerOutput := string(runCommandInContainer(t, "ps", false, "axu"))
		assert.Len(t, strings.Split(containerOutput, "\n"), 4)
		assert.Contains(t, containerOutput, "1 root")
		assert.Contains(t, containerOutput, "root      0:00 ps axu")
	}
}

func hostSeesProcessesInContainerButWithDifferentIDs(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		runCommandInContainer(t, "sleep", true, "infinity")
		containerOutputBytes := runCommandInContainer(t, "ps", false, "axu")

		regExp := regexp.MustCompile(`(?m)^\s*?(\d*?)\s*?root\s*?0:00\s*?sleep infinity$`)

		matches := regExp.FindSubmatch(containerOutputBytes)
		pidInContainer, err := strconv.Atoi(string(matches[1]))
		require.NoError(t, err)

		pidInHost, err := utils.ProcessParentPIDInHostByChildPIDInContainer("sleep", pidInContainer)
		require.NoError(t, err)

		pidStr := strconv.Itoa(pidInHost)
		commandName, err := exec.Command("ps", "-p", pidStr, "-o", "comm=").CombinedOutput()
		require.NoError(t, err)

		commandName = bytes.TrimSpace(commandName)

		assert.Equal(t, []byte("sleep"), commandName)
		assert.NotEqual(t, pidInHost, pidInContainer)
	}
}

func killProcessRunningInContainerFromHost(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		runCommandInContainer(t, "sleep", true, "infinity")
		containerOutputBytes := runCommandInContainer(t, "ps", false, "axu")

		regExp := regexp.MustCompile(`(?m)^\s*?(\d*?)\s*?root\s*?0:00\s*?sleep infinity$`)

		matches := regExp.FindSubmatch(containerOutputBytes)
		pidInContainer, err := strconv.Atoi(string(matches[1]))
		require.NoError(t, err)

		pidInHost, err := utils.ProcessParentPIDInHostByChildPIDInContainer("sleep", pidInContainer)
		require.NoError(t, err)

		pidInHostStr := strconv.Itoa(pidInHost)
		err = exec.Command("kill", "-9", pidInHostStr).Run()
		require.NoError(t, err)

		containerOutputBytes = runCommandInContainer(t, "ps", false, "axu")
		assert.NotContains(t, string(containerOutputBytes), "sleep infinity")
	}
}

func sameHostnameInHostAndInContainer(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		hostnameOnHost, err := os.Hostname()
		require.NoError(t, err)
		domainNameOnHost, err := os.ReadFile("/proc/sys/kernel/domainname")
		require.NoError(t, err)
		hostnameInContainer := string(runCommandInContainer(t, "hostname", false))
		domainNameInContainer := string(runCommandInContainer(t, "cat", false, "/proc/sys/kernel/domainname"))

		assert.NotEqual(t, hostnameOnHost, hostnameInContainer)
		assert.NotEqual(t, string(domainNameOnHost), domainNameInContainer)
		assert.Equal(t, containerHostName, hostnameInContainer)
		assert.Equal(t, containerDomainName, domainNameInContainer)
	}
}

func containerDoesNotSeeQueueOfHost(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		posixQueueBinOnHost := filepath.Join(testRootFS, "bin", "posix-queue")
		queueName := strconv.Itoa(time.Now().Nanosecond())
		cmd := exec.Command(posixQueueBinOnHost, "send-to-queue", queueName, "test data")
		err := cmd.Run()
		require.NoError(t, err)

		assert.Equal(t,
			fmt.Sprintf("queue `/%s` is not existed", queueName),
			string(runCommandInContainer(t, "/bin/posix-queue", false, "is-queue-existed", queueName)),
		)
	}
}

func pingHostFromContainer(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		ip, _, err := net.ParseCIDR(example2.HostIP)
		require.NoError(t, err)
		containerOutput := string(runCommandInContainer(t, "ping", false, "-c", "1", ip.String()))
		assert.Contains(t, containerOutput, "64 bytes from 10.0.0.2")
	}
}

func limitMemoryByCgroup(runCommandInContainer runCommandInContainerFn) asserFn {
	return func(t *testing.T) {
		t.Helper()

		runCommandInContainer(t, "memory-eater", true)

		millisecondsWaitToEatMemory := 50

		time.Sleep(time.Millisecond * time.Duration(millisecondsWaitToEatMemory))

		out, err := exec.Command("systemd-cgtop", "-r", "-m", "-n", "1", "container").CombinedOutput()
		require.NoError(t, err)

		matches := regexp.MustCompile(`(?m)^container/\d*.*?(\d{2,})`).FindAllSubmatch(out, -1)

		for _, match := range matches {
			ramConsumption, err := strconv.Atoi(string(match[1]))
			require.NoError(t, err)

			assert.LessOrEqual(t, ramConsumption, 50000000)
		}
	}
}

func prepareHelpers(t *testing.T, containerRootPath string) {
	t.Helper()

	posixQueueBinaryPath, err := filepath.Abs("../../../utils/os/queue/cmd/posix-queue")
	require.NoError(t, err)

	memoryEaterBinaryPath, err := filepath.Abs("../../../utils/misc/memory_eater/memory-eater")
	require.NoError(t, err)

	err = copyFiles([]string{posixQueueBinaryPath, memoryEaterBinaryPath}, filepath.Join(containerRootPath, "bin"))
	require.NoError(t, err)
}

func command(t *testing.T, cmd string, startInBackground bool, arguments ...string) []byte {
	t.Helper()

	data, err := json.Marshal(example2.Command{
		Cmd:               cmd,
		StartInBackground: startInBackground,
		Arguments:         arguments,
	})
	require.NoError(t, err)

	return example2.PrepareCommand(data)
}

func copyFiles(sources []string, destination string) error {
	sources = append(sources, destination)

	out, err := exec.Command("cp", sources...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("commnad `cp` error: %w %s", err, out)
	}

	return nil
}
