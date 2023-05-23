package example2

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRootFS = "./test_data/rootfs"

type containerOutput struct {
	PID                int      `json:"pid"`
	MountNamespace     int      `json:"mount_namespace"`
	File               string   `json:"file"`
	FileContent        string   `json:"fileContent"`
	TargetDirFilesList []string `json:"targetDirFilesList"`
}

// CHECK HYPOTHESIS OF RUNNING THAT NOT IN TEST AND THAT'S IT

func TestContainer(t *testing.T) {
	rootPath, err := filepath.Abs(testRootFS)
	require.NoError(t, err)

	sourceDir, targetDir, fileWithContentName := "/tmp/source_dir", "/tmp/target_dir", "file_with_content_in_container"

	parentPID := os.Getpid()

	parentSocket, childSocket := test.SockPair(t)
	defer func() {
		err = parentSocket.Close()
		require.NoError(t, err)
	}()

	err = container(
		rootPath,
		childSocket,
		"/bin/mount_in_container/mount_in_container",
		sourceDir,
		targetDir,
		fileWithContentName,
		"content in file",
	)
	require.NoError(t, err)

	tests := []struct {
		name       string
		assertFunc func(t *testing.T, containerOutput containerOutput)
	}{
		{
			name: "host and container are in different mount namespaces",
			assertFunc: func(t *testing.T, containerOutput containerOutput) {
				namespace1, err := utils.MountNamespaceInodeNumber(parentPID)
				require.NoError(t, err)
				assert.NotEqual(t, namespace1, containerOutput.MountNamespace)
			},
		},
		{
			name: "host does not see the mount point in the container",
			assertFunc: func(t *testing.T, containerOutput containerOutput) {
				supposedFilePathOnHost := filepath.Join(testRootFS, targetDir, fileWithContentName)
				_, err = os.Stat(supposedFilePathOnHost)
				assert.True(t, os.IsNotExist(err))
			},
		},
	}

	if parentPID == os.Getpid() {
		err = childSocket.Close()
		require.NoError(t, err)

		data, err := io.ReadAll(parentSocket)
		require.NoError(t, err)

		output := new(containerOutput)

		err = json.Unmarshal(data, output)
		require.NoError(t, err)

		for _, tt := range tests {
			testCase := tt
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				testCase.assertFunc(t, *output)
			})
		}
	}
}
