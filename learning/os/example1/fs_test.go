package example1_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFSSourcePath = "./test_data/test_fs_source"

func TestChroot(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("./test_data/test_bin_source/chroot", "./test_data/test_fs_source")
	output, err := cmd.Output()
	require.NoError(t, err)

	resp := new(struct {
		DirsBeforeChroot []string `json:"dirs_before_chroot"`
		DirsAfterChroot  []string `json:"dirs_after_chroot"`
	})

	err = json.Unmarshal(output, resp)
	require.NoError(t, err)

	assert.Greater(t, len(resp.DirsBeforeChroot), 3)
	assert.Equal(t, 3, len(resp.DirsAfterChroot))
}

/*func TestPivotRoot(t *testing.T) {
	t.Parallel()

	err := syscall.Unshare(syscall.CLONE_NEWNS)
	require.NoError(t, err)

	newRoot := t.TempDir()
	oldRoot := "/old_root_fs"
	oldRoot = filepath.Join(newRoot, oldRoot)
	err = os.MkdirAll(oldRoot, 0o755)
	require.NoError(t, err)

	test.PreventSharedPropagationToRootPoint(t, newRoot)
	test.MakeDirRootPoint(t, newRoot)

	err = syscall.PivotRoot(newRoot, oldRoot)
	require.NoError(t, err)
}*/

func TestMount(t *testing.T) {
	t.Parallel()

	flags := syscall.MS_REC | syscall.MS_BIND | syscall.MS_PRIVATE
	mountPoint := t.TempDir()

	err := syscall.Mount(testFSSourcePath, mountPoint, "tmpfs", uintptr(flags), "mode=0700")
	require.NoError(t, err)

	t.Cleanup(func() {
		err = syscall.Unmount(mountPoint, syscall.MNT_FORCE|syscall.MNT_DETACH)
		require.NoError(t, err)
	})

	var sourceEntries, targetEntries []string

	err = filepath.Walk(testFSSourcePath, func(path string, info fs.FileInfo, err error) error {
		sourceEntries = append(sourceEntries, info.Name())

		return nil
	})
	sourceEntries = sourceEntries[1:]
	err = filepath.Walk(mountPoint, func(path string, info fs.FileInfo, err error) error {
		targetEntries = append(targetEntries, info.Name())

		return nil
	})
	targetEntries = targetEntries[1:]

	assert.ElementsMatch(t, sourceEntries, targetEntries)
}

func TestRootFSFilesWillBeRevertedAfterUnmountTempFS(t *testing.T) {
	t.Parallel()

	mountDir, sourceDir := t.TempDir(), t.TempDir()

	sourceExpectedContentBeforeMount := addTestDirectoriesToDirectory(t, sourceDir, "_source_dir_initial_content")
	sourceActualContentBeforeMount := contentInDirectory(t, sourceDir)
	require.ElementsMatch(t, sourceExpectedContentBeforeMount, sourceActualContentBeforeMount)

	mountPointExpectedContentBeforeMount := addTestDirectoriesToDirectory(t, mountDir, "_mount_point_initial_content")
	mountPointActualContentBeforeMount := contentInDirectory(t, mountDir)
	require.ElementsMatch(t, mountPointExpectedContentBeforeMount, mountPointActualContentBeforeMount)

	test.MountFSToDirectory(t, sourceDir, mountDir)
	mountPointExpectedContentAfterMount := append(
		contentInDirectory(t, sourceDir),
		addTestDirectoriesToDirectory(t, mountDir, "_after_mount")...,
	)
	mountPointActualContentAfterMount := contentInDirectory(t, mountDir)
	require.ElementsMatch(t, mountPointExpectedContentAfterMount, mountPointActualContentAfterMount)

	err := syscall.Unmount(mountDir, syscall.MNT_FORCE|syscall.MNT_DETACH)
	require.NoError(t, err)

	actualContentInDirectoryAfterUnmount := contentInDirectory(t, mountDir)
	sourceActualContentAfterUnmount := contentInDirectory(t, sourceDir)
	assert.ElementsMatch(t, actualContentInDirectoryAfterUnmount, mountPointActualContentBeforeMount)
	assert.ElementsMatch(t, sourceActualContentAfterUnmount, mountPointExpectedContentAfterMount)
}

func addTestDirectoriesToDirectory(t *testing.T, dirPath, testDirsNamePostfix string) []string {
	t.Helper()

	dirsToAdd := [3]string{
		"test_dir_1" + testDirsNamePostfix,
		"test_dir_2" + testDirsNamePostfix,
		"test_dir_n" + testDirsNamePostfix,
	}

	for _, dir := range dirsToAdd {
		err := os.MkdirAll(dirPath+"/"+dir, 0o755)
		require.NoError(t, err)
	}

	return dirsToAdd[:]
}

func contentInDirectory(t *testing.T, dirPath string) []string {
	t.Helper()

	var content []string

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		content = append(content, info.Name())

		return nil
	})
	require.NoError(t, err)

	return content[1:]
}
