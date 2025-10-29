package net_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/net"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVeth(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		err := net.DeleteLink("test_veth")
		require.NoError(t, err)
	})

	veth, err := net.SetupVeth("test_veth", "10.0.0.10/16", "test_peer", "10.0.0.20/16", "testNS1")
	require.NoError(t, err)

	assert.Equal(t, "test_veth", veth.Name)
}
