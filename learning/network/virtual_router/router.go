package virtual_router

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"testing"

	"github.com/phayes/freeport"
	"github.com/r-erema/go_sendbox/pkg/syscall"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

type Router struct {
	arpTable map[[ipLengthBytes]byte]net.HardwareAddr
	connPool map[[macLengthBytes]byte]int
}

func NewRouter() *Router {
	return &Router{
		arpTable: make(map[[4]byte]net.HardwareAddr),
		connPool: make(map[[6]byte]int),
	}
}

func (r *Router) Run(t *testing.T, gatewayPort chan<- uint16) int {
	t.Helper()

	socketIn, err := syscall.SocketFD(unix.AF_INET, unix.SOCK_STREAM)
	require.NoError(t, err)

	portInt, err := freeport.GetFreePort()
	require.NoError(t, err)

	port := cast.ToUint16(portInt)

	err = syscall.Bind(socketIn, "127.0.0.1", port)
	require.NoError(t, err)

	err = syscall.Listen(socketIn)
	require.NoError(t, err)

	gatewayPort <- port

	ipPool := []net.IP{
		net.ParseIP("192.168.0.1").To4(),
		net.ParseIP("192.168.0.2").To4(),
		net.ParseIP("192.168.0.3").To4(),
	}

	var allocatedIP net.IP

	for {
		conn, err := syscall.Accept(socketIn)
		require.NoError(t, err)

		require.NotEmpty(t, ipPool)
		allocatedIP, ipPool = ipPool[0], ipPool[1:]

		hostMacBuf, err := syscall.Read(conn, macLengthBytes)
		require.NoError(t, err)

		hostMac := net.HardwareAddr(hostMacBuf)
		r.connPool[[macLengthBytes]byte(hostMac)] = conn

		err = syscall.Write(conn, allocatedIP)
		require.NoError(t, err)

		r.arpTable[[ipLengthBytes]byte(allocatedIP)] = hostMac

		errCh := make(chan error, 1)
		go func(errCh chan error) {
			errCh <- r.handleConnection(t, conn)
		}(errCh)

		select {
		case err := <-errCh:
			require.NoError(t, err)
		default:
		}
	}
}

func (r *Router) handleConnection(t *testing.T, conn int) error {
	t.Helper()

	frameLenBuf, err := syscall.Read(conn, messageLengthBytes)
	if err != nil {
		return fmt.Errorf("failed to read frame length: %w", err)
	}

	frameBuf, err := syscall.Read(conn, int(frameLenBuf[0]))
	if err != nil {
		return fmt.Errorf("failed to read frame: %w", err)
	}

	frame := new(ExternalEthernetFrame)

	err = gob.NewDecoder(bytes.NewReader(frameBuf)).Decode(frame)
	if err != nil {
		return fmt.Errorf("failed to decode frame: %w", err)
	}

	destinationMAC := r.arpTable[[4]byte(frame.IpPacket.DestinationIP)]

	newFrame, err := buildFrame(
		frame.SourceMAC,
		destinationMAC,
		frame.IpPacket.SourceIP,
		frame.IpPacket.DestinationIP,
		frame.IpPacket.Payload,
	)
	if err != nil {
		return fmt.Errorf("failed to build frame: %w", err)
	}

	destinationConn := r.connPool[[6]byte(destinationMAC)]

	err = syscall.Write(destinationConn, []byte{cast.ToUint8(len(newFrame))})
	if err != nil {
		return fmt.Errorf("error writing frame length: %w", err)
	}

	err = syscall.Write(destinationConn, newFrame)
	if err != nil {
		return fmt.Errorf("error writing frame: %w", err)
	}

	return nil
}
