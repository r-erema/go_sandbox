package queue

/*
#include <mqueue.h>
#include <fcntl.h>

mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
	return mq_open(name, oflag, mode, attr);
}
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	defaultMessageSize = 10000
	mode               = 0o660
)

func Open(name string, flags int) (int, error) {
	var err error

	id, err := C.mq_open4(C.CString(name), C.int(flags), C.int(mode), nil)
	if err != nil {
		return -1, fmt.Errorf("queue opening error: %w", err)
	}

	return int(id), nil
}

func Close(queueDescriptor int) error {
	var err error

	_, err = C.mq_close(C.int(queueDescriptor))
	if err != nil {
		return fmt.Errorf("closing queue(%d) error: %w", queueDescriptor, err)
	}

	return nil
}

func Send(queueID int, data []byte, priority uint) error {
	var err error

	byteStr := *(*string)(unsafe.Pointer(&data))

	_, err = C.mq_send(C.int(queueID), C.CString(byteStr), C.size_t(len(data)), C.uint(priority))
	if err != nil {
		return fmt.Errorf("sending to queue error: %w", err)
	}

	return nil
}

func Receive(queueID int) ([]byte, error) {
	var (
		priority C.uint
		err      error
	)

	buf := (*C.char)(C.malloc(C.size_t(defaultMessageSize)))

	size, err := C.mq_receive(C.int(queueID), buf, C.size_t(defaultMessageSize), &priority)
	if err != nil {
		return nil, fmt.Errorf("receiving from queue error: %w", err)
	}

	return C.GoBytes(unsafe.Pointer(buf), C.int(size)), nil //nolint: nlreturn
}

func Exists(queueName string) bool {
	_, err := Open(queueName, syscall.O_RDONLY)

	return err == nil
}
