package example1_test

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

	err := os.MkdirAll(imageSourcesPath, os.ModePerm)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.RemoveAll(imageSourcesPath)
		require.NoError(t, err)
	})

	tarFile := createTarball(t)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	f, err := os.Open(tarFile.Name())
	require.NoError(t, err)

	_, err = cli.ImageImport(
		context.Background(),
		types.ImageImportSource{Source: f, SourceName: "-"},
		imageImportedName,
		types.ImageImportOptions{},
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = cli.ImageRemove(context.Background(), imageImportedName, types.ImageRemoveOptions{})
		require.NoError(t, err)
	})

	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: imageImportedName,
		Cmd:   []string{"/parrot"},
	}, &container.HostConfig{}, nil, nil, "")

	t.Cleanup(func() {
		timeout := time.Duration(0)
		err = cli.ContainerStop(context.Background(), cont.ID, &timeout)
		require.NoError(t, err)
		err = cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{})
		require.NoError(t, err)
	})

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
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
