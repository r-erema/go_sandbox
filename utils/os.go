package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MountNamespaceInodeNumber(pid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/ns/mnt", pid)

	link, err := os.Readlink(path)
	if err != nil {
		return -1, fmt.Errorf("reading link `%s` error: %w", path, err)
	}

	IDStr := strings.Replace(link, "mnt:[", "", -1)
	IDStr = strings.Replace(IDStr, "]", "", -1)

	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		return -1, fmt.Errorf("string to integer conversion error: %w", err)
	}

	return ID, nil
}

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

func FilesInDir(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reding dir error: %w", err)
	}

	files := make([]string, len(entries))

	for i := range entries {
		files[i] = entries[i].Name()
	}

	return files, nil
}
