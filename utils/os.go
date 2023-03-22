package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func MountNamespaceInodeNumber(pid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/ns/mnt", pid)

	link, err := os.Readlink(path)
	if err != nil {
		return -1, fmt.Errorf("reading link `%s` error: %w", path, err)
	}

	ID, err := ParseNSID("mnt:[", link)
	if err != nil {
		return -1, fmt.Errorf("parsing link ID `%s` error: %w", link, err)
	}

	return ID, nil
}

func NetworkNamespaceInodeNumber(pid, tid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid)

	link, err := os.Readlink(path)
	if err != nil {
		return -1, fmt.Errorf("reading link `%s` error: %w", path, err)
	}

	ID, err := ParseNSID("net:[", link)
	if err != nil {
		return -1, fmt.Errorf("parsing link ID `%s` error: %w", link, err)
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

func ParseNSID(prefix, link string) (int, error) {
	IDStr := strings.Replace(link, prefix, "", -1)
	IDStr = strings.Replace(IDStr, "]", "", -1)

	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		return -1, fmt.Errorf("string to integer conversion error: %w", err)
	}

	return ID, nil
}

func GrepPIDInHostPIDNSAndChildPIDNS(command string) (map[int]int, error) {
	cmd := exec.Command("pgrep", command)

	pgrepOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command pgrep output error: %w", err)
	}

	pgrepOutput = bytes.TrimSpace(pgrepOutput)

	regExp := regexp.MustCompile(`(?m)NSpid:\s*?(\d*?)\s*?(\d*)$`)

	parentChildPIDsMap := make(map[int]int)

	for _, pidBytes := range bytes.Split(pgrepOutput, []byte{'\n'}) {
		data, err := os.ReadFile(fmt.Sprintf("/proc/%s/status", string(pidBytes)))
		if err != nil {
			return nil, fmt.Errorf("reading file error: %w", err)
		}

		matches := regExp.FindSubmatch(data)

		parentPID, err := strconv.Atoi(string(matches[1]))
		if err != nil {
			return nil, fmt.Errorf("parent PID string to integer conversion error: %w", err)
		}

		childPID, err := strconv.Atoi(string(matches[2]))
		if err != nil {
			return nil, fmt.Errorf("child PID string to integer conversion error: %w", err)
		}

		parentChildPIDsMap[parentPID] = childPID
	}

	return parentChildPIDsMap, nil
}

var errPIDNotFound = errors.New("PID not found")

func ProcessParentPIDInHostByChildPIDInContainer(processName string, childPIDInContainer int) (int, error) {
	parentChildPIDsMap, err := GrepPIDInHostPIDNSAndChildPIDNS(processName)
	if err != nil {
		return -1, fmt.Errorf("getting parent-child PIDs map for process `%s` error: %w", processName, err)
	}

	var foundParentPID, chPID int
	for foundParentPID, chPID = range parentChildPIDsMap {
		if chPID == childPIDInContainer {
			return foundParentPID, nil
		}
	}

	return -1, errPIDNotFound
}
