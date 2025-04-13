package socket

import (
	"fmt"

	"github.com/r-erema/go_sendbox/pkg/tls"
	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"golang.org/x/sys/unix"
)

func Server(
	ipAddr string,
	port uint16,
	publicKey [32]byte,
	handler func(incomingBuf []byte) []byte,
) error {
	socketFD, err := syscall.SocketFD(unix.AF_INET, unix.SOCK_STREAM)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	if err = syscall.Bind(socketFD, ipAddr, port); err != nil {
		return fmt.Errorf("failed to bind socket fd `%d` to IP `%s`: %w", socketFD, ipAddr, err)
	}

	if err = syscall.Listen(socketFD); err != nil {
		return fmt.Errorf("failed to listen on socket fd `%d`: %w", socketFD, err)
	}

	for {
		connFD, err := syscall.Accept(socketFD)
		if err != nil {
			return fmt.Errorf("failed to accept socket fd `%d`: %w", socketFD, err)
		}

		if err = tls.HandshakeServerSide(connFD, publicKey); err != nil {
			return fmt.Errorf("failed to handshake fd `%d`: %w", connFD, err)
		}

		buf, err := syscall.Read(connFD)
		if err != nil {
			return fmt.Errorf(
				"failed to read connection fd `%d` of socket fd `%d`: %w",
				connFD,
				socketFD,
				err,
			)
		}

		buf = handler(buf)

		if err = syscall.Write(connFD, buf); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
}
