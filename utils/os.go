package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func NetworkNamespaceInodeNumber(pid, tid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid)

	link, err := os.Readlink(path)
	if err != nil {
		return -1, fmt.Errorf("reading link `%s` error: %w", path, err)
	}

	IDStr := strings.Replace(link, "net:[", "", -1)
	IDStr = strings.Replace(IDStr, "]", "", -1)

	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		return -1, fmt.Errorf("string to integer conversion error: %w", err)
	}

	return ID, nil
}

func NewNetworkNamespaceDescriptor(pid, tid int) (uintptr, error) {
	path := fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid)

	file, err := os.Open(path)
	if err != nil {
		return uintptr(0), fmt.Errorf("opening file `%s` error: %w", path, err)
	}

	return file.Fd(), nil
}
