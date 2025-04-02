package socket

import (
	"fmt"

	"github.com/r-erema/go_sendbox/pkg/tls"
	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"golang.org/x/sys/unix"
)

func ClientCall(ipAddr string, port uint16, publicKey [32]byte, message []byte) ([]byte, error) {
	socketFD, err := syscall.SocketFD(unix.AF_INET, unix.SOCK_STREAM)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %w", err)
	}

	if err = syscall.Connect(socketFD, ipAddr, port); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err = tls.HandshakeClientSide(socketFD, []string{"localhost.dev"}, publicKey); err != nil {
		return nil, fmt.Errorf("tls handshake failure: %w", err)
	}

	if err = syscall.Write(socketFD, message); err != nil {
		return nil, fmt.Errorf("failed to write to file: %w", err)
	}

	buf, err := syscall.Read(socketFD)
	if err != nil {
		return nil, fmt.Errorf("failed to read connection fd `%d` of socket fd `%d`: %w", socketFD, socketFD, err)
	}

	return buf, nil
}
