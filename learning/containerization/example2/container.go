package example2

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/r-erema/go_sendbox/utils"
)

func container(rootPath string, commandOutput *os.File, command string, args ...string) error {
	pid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		panic(fmt.Sprintf("syscall fork error: %d", errno))
	}

	childProcess := pid == 0

	if childProcess {
		if err := setupFileSystem(rootPath); err != nil {
			return fmt.Errorf("setup file system error: %s", err)
		}

		if err := setHostNameAndDomain(); err != nil {
			return fmt.Errorf("set hostname and domain error: %s", err)
		}

		log.Printf("[Container] PID: %d", os.Getpid())

		log.Printf("[Container] Running command: %s", command)
		if err := runContainer(commandOutput.Fd(), command, args...); err != nil {
			return fmt.Errorf("run runContainer error: %s", err)
		}

		ns, err := utils.MountNamespaceInodeNumber(os.Getpid())
		if err != nil {
			return fmt.Errorf("getting mount namespace inode number: %s", err)
		}
		log.Printf("[Container] Mount namespace: %d", ns)
	}

	return nil
}

func setHostNameAndDomain() error {
	if err := syscall.Sethostname([]byte("container_hostname")); err != nil {
		return fmt.Errorf("set hostname error: %s", err)
	}

	if err := syscall.Setdomainname([]byte("domain_hostname")); err != nil {
		return fmt.Errorf("set domainname error: %s", err)
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
		return fmt.Errorf("error: %w, unable to change the propagation type of all mount points under `/` to MS_SLAVE"+
			"in order prevent container propagate its changes to host, "+
			"but changes from the host can be received", err)
	}

	return nil
}

func changePropagationTypeToShared() error {
	err := syscall.Mount("", "/", "", syscall.MS_SHARED|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("error: %w, unable to change the propagation type of of all mount points under `/` to MS_SHARED"+
			"in order propagate container's to children if any", err)
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
	err := syscall.Mount("proc", "/proc", "proc", syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV, "")
	if err != nil {
		return fmt.Errorf("mount MS_NOSUID | MS_NOEXEC | MS_NODEV error: %w", err)
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
		return fmt.Errorf("mount proc fs error: %s", err)
	}

	return nil
}

func runContainer(outputSocketFD uintptr, command string, args ...string) error {
	args = append([]string{command}, args...)
	args = append(args, strconv.Itoa(int(outputSocketFD)))

	err := syscall.Exec(command, args, os.Environ())
	if err != nil {
		return fmt.Errorf("exec command `%s` error: %w", command, err)
	}
	return nil
}
