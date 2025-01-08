package syscall

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

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
