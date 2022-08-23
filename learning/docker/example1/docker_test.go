package example1_test

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_docker(t *testing.T) {
	t.Parallel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	ping, err := cli.Ping(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, ping)
}
