package serveronsyscalls

import (
	"bytes"
	"fmt"

	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"golang.org/x/sys/unix"
)

func Server(ipAddr string, port uint16, incomingData chan<- []byte) error {
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

	connFD, err := syscall.Accept(socketFD)
	if err != nil {
		return fmt.Errorf("failed to accept socket fd `%d`: %w", socketFD, err)
	}

	buf, err := syscall.Read(connFD)
	if err != nil {
		return fmt.Errorf("failed to read connection fd `%d` of socket fd `%d`: %w", connFD, socketFD, err)
	}
	incomingData <- bytes.Trim(buf, "\x00")

	return nil
}
