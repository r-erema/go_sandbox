package example2

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

const bufferSize = 100

func Server(ipString string, port int, incomingData chan<- []byte) error {
	ipAddr := net.ParseIP(ipString)

	fileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0) //nolint: nosnakecase
	if err != nil {
		return fmt.Errorf("syscall socket error: %w", os.NewSyscallError("socket", err))
	}

	defer func() {
		err = syscall.Close(fileDescriptor)
		if err != nil {
			log.Println(os.NewSyscallError("close", err))
		}
	}()

	if err = syscall.SetsockoptInt(fileDescriptor, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil { //nolint: nosnakecase
		return fmt.Errorf("syscall setsockopt error: %w", os.NewSyscallError("setsockopt", err))
	}

	sa := &syscall.SockaddrInet4{Port: port, Addr: [4]byte{}}
	copy(sa.Addr[:], ipAddr)

	if err = syscall.Bind(fileDescriptor, sa); err != nil {
		return fmt.Errorf("syscall bind error: %w", os.NewSyscallError("bind", err))
	}

	if err = syscall.Listen(fileDescriptor, syscall.SOMAXCONN); err != nil {
		return fmt.Errorf("syscall listen error: %w", os.NewSyscallError("listen", err))
	}

	nfd, _, err := syscall.Accept(fileDescriptor)
	if err != nil {
		return fmt.Errorf("syscall accept error: %w", os.NewSyscallError("accept", err))
	}

	defer func() {
		err = syscall.Close(nfd)
		if err != nil {
			log.Println(os.NewSyscallError("close", err))
		}
	}()

	buffer := make([]byte, bufferSize)

	_, err = syscall.Read(nfd, buffer)
	if err != nil {
		return fmt.Errorf("syscall read error: %w", os.NewSyscallError("read", err))
	}

	incomingData <- bytes.Trim(buffer, "\x00")

	return nil
}
