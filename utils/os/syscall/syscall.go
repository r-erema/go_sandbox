package syscall

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const defaultBufferSize = 4096

func SocketFD(domain, sockType int) (int, error) {
	switch domain {
	case unix.AF_INET:
	case unix.AF_INET6:
	case unix.AF_UNIX:
	default:
		return -1, fmt.Errorf("unsupported domain `%d`: %w", domain, os.NewSyscallError("SYS_SOCKET", nil))
	}

	switch sockType {
	case unix.SOCK_STREAM:
	case unix.SOCK_DGRAM:
	default:
		return -1, fmt.Errorf("unsupported socket type: `%d`: %w", sockType, os.NewSyscallError("SYS_SOCKET", nil))
	}

	fd, _, errno := unix.Syscall(unix.SYS_SOCKET, uintptr(domain), uintptr(sockType), uintptr(0))
	if errno != 0 {
		return -1, fmt.Errorf("syscall SYS_SOCKET error: %w", errno)
	}

	return int(fd), nil
}

func Bind(socketFD int, ipAddr string, port uint16) error {
	ip := net.ParseIP(ipAddr)

	raw := &syscall.RawSockaddrInet4{
		Family: unix.AF_INET,
		Port:   port,
		Addr:   [4]byte(ip),
	}

	p := (*[2]byte)(unsafe.Pointer(&raw.Port))
	offset := 8
	p[0] = byte(port >> offset)
	p[1] = byte(port)

	_, _, errno := unix.Syscall(unix.SYS_BIND, uintptr(socketFD), uintptr(unsafe.Pointer(raw)), uintptr(syscall.SizeofSockaddrInet4))
	if errno != 0 {
		return fmt.Errorf("syscall SYS_BIND error: %w", errno)
	}

	return nil
}

func Listen(socketFD int) error {
	_, _, errno := unix.Syscall(unix.SYS_LISTEN, uintptr(socketFD), uintptr(syscall.SOMAXCONN), 0)
	if errno != 0 {
		return fmt.Errorf("syscall SYS_LISTEN error: %w", errno)
	}

	return nil
}

func Accept(socketFD int) (int, error) {
	rsa := unsafe.Pointer(&unix.RawSockaddrAny{})

	var addrLen uint32 = unix.SizeofSockaddrAny

	fd, _, errno := unix.Syscall(unix.SYS_ACCEPT, uintptr(socketFD), uintptr(rsa), uintptr(unsafe.Pointer(&addrLen)))
	if errno != 0 {
		return -1, fmt.Errorf("syscall SYS_ACCEPT error: %w", errno)
	}

	return int(fd), nil
}

func Read(fd int) ([]byte, error) {
	buf := make([]byte, defaultBufferSize)
	bufStartPtr := unsafe.Pointer(&buf[0])

	_, _, errno := unix.Syscall(unix.SYS_READ, uintptr(fd), uintptr(bufStartPtr), defaultBufferSize)
	if errno != 0 {
		return nil, fmt.Errorf("syscall SYS_READ error for FD %d: %w", fd, errno)
	}

	return buf, nil
}

func ReadFromFileFD(fd int) (string, error) {
	fileLen, _, errno := unix.Syscall(unix.SYS_LSEEK, uintptr(fd), uintptr(0), uintptr(unix.SEEK_END))
	if errno != 0 {
		return "", fmt.Errorf("syscall LSEEK(SEEK_END) error for FD %d: %w", fd, errno)
	}

	_, _, errno = unix.Syscall(unix.SYS_LSEEK, uintptr(fd), uintptr(0), uintptr(unix.SEEK_SET))
	if errno != 0 {
		return "", fmt.Errorf("syscall LSEEK(SEEK_SET) error for FD %d: %w", fd, errno)
	}

	buf := make([]byte, fileLen)
	bufStartPtr := unsafe.Pointer(&buf[0])

	_, _, errno = unix.Syscall(unix.SYS_READ, uintptr(fd), uintptr(bufStartPtr), fileLen)
	if errno != 0 {
		return "", fmt.Errorf("syscall SYS_READ error for FD %d: %w", fd, errno)
	}

	return string(buf), nil
}

func WriteToFileFD(fd int, data []byte) error {
	dataLen := len(data)
	bufStartPtr := unsafe.Pointer(&data[0])

	_, _, errno := unix.Syscall(unix.SYS_WRITE, uintptr(fd), uintptr(bufStartPtr), uintptr(dataLen))
	if errno != 0 {
		return fmt.Errorf("syscall SYS_WRITE error for FD %d: %w", fd, errno)
	}

	return nil
}
