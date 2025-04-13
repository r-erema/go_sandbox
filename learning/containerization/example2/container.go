package example2

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/r-erema/go_sendbox/utils"
	utilsNet "github.com/r-erema/go_sendbox/utils/os/net"
	"golang.org/x/sys/unix"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	StreamDelimiter             = '|'
	bridgeName                  = "host_br"
	hostEthName                 = "host_eth"
	containerEthName            = "container_eth"
	bridgeIP                    = "10.0.0.1/8"
	HostIP                      = "10.0.0.2/8"
	containerIP                 = "10.0.0.3/8"
	cgroupPath                  = "/sys/fs/cgroup/container/"
	cgroupMemoryLimitBytes      = 50000000 // 50 MB
	cgroupPathPermissions       = 0o600
	cgroupLowMemoryLimitPortion = 0.75
)

type Command struct {
	Cmd               string   `json:"cmd"`
	StartInBackground bool     `json:"start_in_background"`
	Arguments         []string `json:"arguments"`
}

func PrepareCommand(cmd []byte) []byte {
	return append(cmd, StreamDelimiter)
}

func TrimOutput(output []byte) []byte {
	return bytes.TrimSpace(bytes.Trim(output, string(StreamDelimiter)))
}

func Container(rootPath, hostName, domainName string, commandOutput *os.File) (int, error) {
	waitingHostContainerPipe, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil {
		return -1, fmt.Errorf("creation sockpair error: %w", err)
	}

	pid, _, errno := syscall.Syscall(syscall.SYS_CLONE, uintptr(syscall.SIGCHLD|
		syscall.CLONE_NEWNS|
		syscall.CLONE_NEWPID|
		syscall.CLONE_NEWUTS|
		syscall.CLONE_NEWIPC|
		syscall.CLONE_NEWNET,
	), 0, 0)
	if errno != 0 {
		panic(fmt.Sprintf("syscall fork error: %d", errno))
	}

	childProcess := pid == 0

	if childProcess {
		if err = containerBranch(rootPath, hostName, domainName, commandOutput, waitingHostContainerPipe); err != nil {
			return -1, fmt.Errorf("container branch running error: %w", err)
		}
	} else {
		if err = hostBranch(pid, waitingHostContainerPipe); err != nil {
			return -1, fmt.Errorf("container branch running error: %w", err)
		}
	}

	return int(pid), nil
}

func containerBranch(
	rootPath, hostName, domainName string,
	commandOutput *os.File,
	waitingHostContainerPipe [2]int,
) error {
	log.Printf("[Container] Waiting for the host preparation...")

	if _, err := os.NewFile(uintptr(waitingHostContainerPipe[1]), "").Read(make([]byte, 1)); err != nil {
		return fmt.Errorf("reading from the waitingHostContainerPipe error: %w", err)
	}

	log.Printf("[Container] Host is ready, prceeding...")

	if err := setupNetworkInContainer(); err != nil {
		return fmt.Errorf("setup network error: %w", err)
	}

	if err := setupFileSystem(rootPath); err != nil {
		return fmt.Errorf("setup file system error: %w", err)
	}

	if err := setHostAndDomainName(hostName, domainName); err != nil {
		return fmt.Errorf("set hostname and domain error: %w", err)
	}

	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("reading hostname error: %w", err)
	}

	log.Printf("[Container] Hostname: %s", host)

	domain, err := os.ReadFile("/proc/sys/kernel/domainname")
	if err != nil {
		return fmt.Errorf("reading hostname error: %w", err)
	}

	log.Printf("[Container] Domainname: %s", domain)
	log.Printf("[Container] PID: %d", os.Getpid())

	namespace, err := utils.MountNamespaceInodeNumber(os.Getpid())
	if err != nil {
		return fmt.Errorf("getting mount namespace inode number: %w", err)
	}

	log.Printf("[Container] Mount namespace: %d", namespace)

	namespace, err = utils.NetworkNamespaceInodeNumber(os.Getpid(), unix.Gettid())
	if err != nil {
		return fmt.Errorf("getting network namespace inode number: %w", err)
	}

	log.Printf("[Container] Network namespace: %d", namespace)

	runProcess(commandOutput.Fd())

	return nil
}

func hostBranch(containerPID uintptr, waitingHostContainerPipe [2]int) error {
	log.Printf("[Host] PID: %d", os.Getpid())

	namespace, err := utils.MountNamespaceInodeNumber(os.Getpid())
	if err != nil {
		return fmt.Errorf("getting mount namespace inode number error: %w", err)
	}

	log.Printf("[Host] Mount namespace: %d", namespace)

	namespace, err = utils.NetworkNamespaceInodeNumber(os.Getpid(), unix.Gettid())
	if err != nil {
		return fmt.Errorf("getting network namespace inode number error: %w", err)
	}

	log.Printf("[Host] Network namespace: %d", namespace)

	if err = setupNetworkOnHost(containerPID); err != nil {
		return fmt.Errorf("setup network error: %w", err)
	}

	if err = setupCgroup(int(containerPID), cgroupMemoryLimitBytes); err != nil {
		return fmt.Errorf("setup cgroup error: %w", err)
	}

	log.Printf("[Host] Preparation is done, sending signal to the container")

	if _, err = os.NewFile(uintptr(waitingHostContainerPipe[0]), "").Write([]byte{0}); err != nil {
		return fmt.Errorf("writing to the waitingHostContainerPipe error: %w", err)
	}

	return nil
}

func KillContainer(containerPID int) error {
	containerProcess, err := os.FindProcess(containerPID)
	if err != nil {
		return fmt.Errorf("lookup container process error: %w", err)
	}

	if err = containerProcess.Kill(); err != nil {
		return fmt.Errorf("killing container process error: %w", err)
	}

	if err = utilsNet.RemoveBridge(bridgeName); err != nil {
		return fmt.Errorf("bridge setup error: %w", err)
	}

	if err = removeCgroup(cgroupPathOfProcess(containerPID)); err != nil {
		return fmt.Errorf("removal container cgroup error: %w", err)
	}

	return nil
}

func setHostAndDomainName(hostName, domainName string) error {
	if err := syscall.Sethostname([]byte(hostName)); err != nil {
		return fmt.Errorf("set hostname error: %w", err)
	}

	if err := syscall.Setdomainname([]byte(domainName)); err != nil {
		return fmt.Errorf("set domainname error: %w", err)
	}

	return nil
}

func createNewNamespace() error {
	err := syscall.Unshare(syscall.CLONE_NEWNS)
	if err != nil {
		return fmt.Errorf("cloning namespace error: %w", err)
	}

	return nil
}

func changePropagationTypeToSlave() error {
	err := syscall.Mount("", "/", "", syscall.MS_SLAVE|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf(
			"error: %w, unable to change the propagation type of all mount points under `/` to MS_SLAVE"+
				"in order prevent container propagate its changes to host, "+
				"but changes from the host can be received",
			err,
		)
	}

	return nil
}

func changePropagationTypeToShared() error {
	err := syscall.Mount("", "/", "", syscall.MS_SHARED|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf(
			"error: %w, unable to change the propagation type of all mount points under `/` to MS_SHARED"+
				"in order propagate container's to children if any",
			err,
		)
	}

	return nil
}

func bindDirToItself(dir string) error {
	err := syscall.Mount(dir, dir, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("mount MS_BIND dir `%s` error: %w", dir, err)
	}

	return nil
}

func moveMountPoint(source, destination string) error {
	err := syscall.Mount(source, destination, "", syscall.MS_MOVE, "")
	if err != nil {
		return fmt.Errorf("mount MS_MOVE from `%s` to `%s` error: %w", source, destination, err)
	}

	return nil
}

func mountProcFS() error {
	err := syscall.Mount(
		"proc",
		"/proc",
		"proc",
		syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV,
		"",
	)
	if err != nil {
		return fmt.Errorf("mount MS_NOSUID | MS_NOEXEC | MS_NODEV error: %w", err)
	}

	return nil
}

func setupNetworkOnHost(containerPID uintptr) error {
	err := utilsNet.SetupBridge(bridgeName, bridgeIP)
	if err != nil {
		return fmt.Errorf("bridge setup error: %w", err)
	}

	_, err = utilsNet.SetupVeth(
		hostEthName,
		HostIP,
		containerEthName,
		strconv.Itoa(int(containerPID)),
	)
	if err != nil {
		return fmt.Errorf("veth setup error: %w", err)
	}

	err = utilsNet.AttachDeviceToBridge(hostEthName, bridgeName)
	if err != nil {
		return fmt.Errorf("veth attach to the bridge error: %w", err)
	}

	return nil
}

func setupNetworkInContainer() error {
	_, err := utilsNet.SetupLoopBackInterface()
	if err != nil {
		return fmt.Errorf("loopback interface setup error: %w", err)
	}

	if err = utilsNet.AddIPAddrToInterface(containerIP, containerEthName); err != nil {
		return fmt.Errorf(
			"adding ip `%s` address to the interface `%s` error: %w",
			containerIP,
			containerEthName,
			err,
		)
	}

	if err = utilsNet.EnableDevice(containerEthName); err != nil {
		return fmt.Errorf("enabling device `%s` error: %w", containerEthName, err)
	}

	ip, _, err := net.ParseCIDR(bridgeIP)
	if err != nil {
		return fmt.Errorf("parse CIDR `%s` error: %w", bridgeIP, err)
	}

	if err = utilsNet.SetDefaultGateway(ip.String()); err != nil {
		return fmt.Errorf("set default gateway error: %w", err)
	}

	return nil
}

func setupFileSystem(rootFS string) error {
	err := createNewNamespace()
	if err != nil {
		return fmt.Errorf("cloning namespace error: %w", err)
	}

	err = changePropagationTypeToSlave()
	if err != nil {
		return fmt.Errorf("changing propagation type to SLAVE error: %w", err)
	}

	err = bindDirToItself(rootFS)
	if err != nil {
		return fmt.Errorf("binding directory error: %w", err)
	}

	err = os.Chdir(rootFS)
	if err != nil {
		return fmt.Errorf("chroot error: %w", err)
	}

	err = moveMountPoint(rootFS, "/")
	if err != nil {
		return fmt.Errorf("moving mount point to root error: %w", err)
	}

	err = syscall.Chroot(".")
	if err != nil {
		return fmt.Errorf("chroot error: %w", err)
	}

	err = os.Chdir("/")
	if err != nil {
		return fmt.Errorf("chroot error: %w", err)
	}

	err = changePropagationTypeToShared()
	if err != nil {
		return fmt.Errorf("changing propagation type to SHARED error: %w", err)
	}

	err = mountProcFS()
	if err != nil {
		return fmt.Errorf("mount proc fs error: %w", err)
	}

	return nil
}

func runProcess(outputSocketFD uintptr) {
	sock := os.NewFile(outputSocketFD, "")

	log.Printf("[Container] Process started")

	for {
		r := bufio.NewReader(sock)

		data, err := r.ReadBytes(StreamDelimiter)
		if err != nil {
			log.Print(fmt.Errorf("[Container] Reading bytes error: %w", err))
		}

		command := new(Command)

		err = json.Unmarshal(TrimOutput(data), command)
		if err != nil {
			log.Print(fmt.Errorf("[Container] Unmarshalling command error: %w", err))
		}

		runExternalCommand(command, sock)
	}
}

func runExternalCommand(command *Command, sock *os.File) {
	commandName, args := command.Cmd, command.Arguments
	cmd := exec.Command(commandName, args...) //nolint:gosec // disable G115

	log.Printf("[Container] Running command: %s with arguments: %s", commandName, args)

	if command.StartInBackground {
		err := cmd.Start()
		if err != nil {
			log.Print(
				fmt.Errorf(
					"[Container] Command `%v` starting error: %w(%s)",
					command,
					err,
					cmd.Stderr,
				),
			)
		}

		_, err = sock.Write([]byte{StreamDelimiter})
		if err != nil {
			log.Print(
				fmt.Errorf(
					"[Container] Command `%v` output to sock writing error: %w",
					command,
					err,
				),
			)
		}

		return
	}

	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(fmt.Errorf("[Container] Command `%v` running error: %w(%s)", command, err, data))
	}

	log.Printf("[Container] Command output: %s", data)

	_, err = sock.Write(append(data, StreamDelimiter))
	if err != nil {
		log.Print(fmt.Errorf("[Container] Writing to output error: %w", err))
	}
}

func cgroupPathOfProcess(pid int) string {
	return cgroupPath + strconv.Itoa(pid)
}

func removeCgroup(cgroupPath string) error {
	if err := os.Remove(cgroupPath); err != nil {
		return fmt.Errorf("removal path `%s` error: %w", cgroupPath, err)
	}

	return nil
}

func setupCgroup(pid, limitRAMBytes int) error {
	cgPath := cgroupPathOfProcess(pid)

	if err := os.MkdirAll(cgPath, cgroupPathPermissions); err != nil {
		return fmt.Errorf("creation cgroup path `%s` error: %w", cgPath, err)
	}

	subtreeControlFile := cgroupPath + "cgroup.subtree_control"
	if err := os.WriteFile(
		subtreeControlFile,
		[]byte("+memory +cpu +cpuset +io +pids"), cgroupPathPermissions,
	); err != nil {
		return fmt.Errorf("enabling memory cgroup controller error: %w", err)
	}

	if limitRAMBytes > 0 {
		memoryLowFile := cgPath + "/memory.low"
		memoryLowVal := []byte(
			strconv.FormatFloat(float64(limitRAMBytes)*cgroupLowMemoryLimitPortion, 'f', -1, 64),
		)

		if err := os.WriteFile(memoryLowFile, memoryLowVal, cgroupPathPermissions); err != nil {
			return fmt.Errorf("writing to the file `%s` error: %w", memoryLowFile, err)
		}

		memoryMaxFile := cgPath + "/memory.max"
		if err := os.WriteFile(memoryMaxFile, []byte(strconv.Itoa(limitRAMBytes)), cgroupPathPermissions); err != nil {
			return fmt.Errorf("writing to the file `%s` error: %w", memoryMaxFile, err)
		}
	}

	cgProcsFile := cgPath + "/cgroup.procs"
	if err := os.WriteFile(cgProcsFile, []byte(strconv.Itoa(pid)), cgroupPathPermissions); err != nil {
		return fmt.Errorf("writing to the file `%s` error: %w", cgProcsFile, err)
	}

	return nil
}
