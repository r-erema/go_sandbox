package example1_test

import (
	"archive/tar"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

const (
	imageSourcesPath  = "/tmp/test_image/layer"
	imageOutputPath   = "/tmp/test_image/image.tar"
	imageImportedName = "parrot"
	binPath           = "./parrot"
)

func TestBuildImage(t *testing.T) {
	t.Parallel()

	err := os.MkdirAll(imageSourcesPath, 0o750)
	require.NoError(t, err)

	defer func() {
		err = os.RemoveAll(imageSourcesPath)
		require.NoError(t, err)
	}()

	tarFile := createTarball(t)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	file, err := os.Open(tarFile.Name())
	require.NoError(t, err)

	ctx := t.Context()

	_, err = cli.ImageImport(
		ctx,
		image.ImportSource{Source: file, SourceName: "-"},
		imageImportedName,
		image.ImportOptions{},
	)
	require.NoError(t, err)

	defer func() {
		_, err = cli.ImageRemove(t.Context(), imageImportedName, image.RemoveOptions{})
		require.NoError(t, err)
	}()

	cont, err := cli.ContainerCreate(t.Context(), &container.Config{
		Image: imageImportedName,
		Cmd:   []string{"/parrot"},
	}, &container.HostConfig{}, nil, nil, "")

	defer func() {
		timeout := 100
		err = cli.ContainerStop(t.Context(), cont.ID, container.StopOptions{Timeout: &timeout})
		require.NoError(t, err)
		err = cli.ContainerRemove(t.Context(), cont.ID, container.RemoveOptions{})
		require.NoError(t, err)
	}()

	err = cli.ContainerStart(t.Context(), cont.ID, container.StartOptions{})
	require.NoError(t, err)
}

func createTarball(t *testing.T) *os.File {
	t.Helper()

	tarFile, err := os.Create(imageOutputPath)
	require.NoError(t, err)

	tarWriter := tar.NewWriter(tarFile)

	binFile, err := os.Open(binPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = binFile.Close()
		require.NoError(t, err)
		err = os.RemoveAll(imageOutputPath)
		require.NoError(t, err)
	})

	stat, err := binFile.Stat()
	require.NoError(t, err)

	header := &tar.Header{
		Name:    "parrot",
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	require.NoError(t, err)

	_, err = io.Copy(tarWriter, binFile)
	require.NoError(t, err)

	return tarFile
}
