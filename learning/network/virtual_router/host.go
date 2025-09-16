package virtual_router

import (
	"bytes"
	"encoding/gob"
	"net"
	"testing"

	"github.com/r-erema/go_sendbox/pkg/syscall"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

type Host struct {
	macAddr net.HardwareAddr
	ipAddr  net.IP
	conn    int
	t       *testing.T
}

func NewHost(t *testing.T, macAddr net.HardwareAddr) *Host {
	t.Helper()

	return &Host{
		macAddr: macAddr,
		t:       t,
	}
}

func (h *Host) IpAddr() net.IP {
	return h.ipAddr
}

func (h *Host) ConnectToGateway(gatewayPort uint16) {
	var err error

	h.conn, err = syscall.SocketFD(unix.AF_INET, unix.SOCK_STREAM)
	require.NoError(h.t, err)

	// connect ot the Router
	err = syscall.Connect(h.conn, "127.0.0.1", gatewayPort)
	require.NoError(h.t, err)

	// Send to the route MAC address of the Host to be registered
	err = syscall.Write(h.conn, h.macAddr)
	require.NoError(h.t, err)

	// getting IP from Router's DHCP
	offeredIP, err := syscall.Read(h.conn, ipLengthBytes)
	require.NoError(h.t, err)

	h.ipAddr = offeredIP
}

func (h *Host) Send(destinationIP net.IP, message []byte) {
	frame, err := buildFrame(h.macAddr, nil, h.ipAddr, destinationIP, message)
	require.NoError(h.t, err)

	err = syscall.Write(h.conn, []byte{byte(len(frame))})
	require.NoError(h.t, err)

	err = syscall.Write(h.conn, frame)
	require.NoError(h.t, err)
}

func (h *Host) Receive() []byte {
	buf, err := syscall.Read(h.conn, messageLengthBytes)
	require.NoError(h.t, err)

	buf, err = syscall.Read(h.conn, int(buf[0]))
	require.NoError(h.t, err)

	frame := new(ExternalEthernetFrame)
	err = gob.NewDecoder(bytes.NewReader(buf)).Decode(frame)
	require.NoError(h.t, err)

	return frame.IpPacket.Payload
}
