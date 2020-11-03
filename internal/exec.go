package internal

import (
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// Exec runs a command inside an existing container.
//
// It asks for container digest, command and arg to run, and a
// detach bool. Container digest can be a prefix of digest.
//
// If detach is true, Exec never wait for command to get done and returns.
func Exec(ctrDigest string, args []string, detach bool) error {
	ctr, err := container.GetContainerByDigest(ctrDigest)
	if err != nil {
		return err
	}
	if len(ctr.Pids) == 0 || ctr.Pids[0] == 0 {
		return errors.Errorf("container %s is not running", ctr.Digest)
	}

	err = setNamespace(ctr.Pids[0], syscall.CLONE_NEWUTS|syscall.CLONE_NEWIPC|syscall.CLONE_NEWPID|syscall.CLONE_NEWNET)
	if err != nil {
		return err
	}

	if err := changeRoot(ctr.RootFS, ctr.Config.WorkingDir); err != nil {
		return err
	}

	mountPoints := []filesystem.MountPoint{
		{Source: "proc", Target: "proc", Type: "proc"},
		{Source: "sysfs", Target: "sys", Type: "sysfs"},
	}
	unmounter, err := filesystem.Mount(mountPoints...)
	if err != nil {
		return err
	}
	defer unmounter()

	newCmd := exec.Command(args[0], args[1:]...)
	newCmd.Stdin, newCmd.Stdout, newCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	newCmd.Env = ctr.Config.Env
	if err := runCommand(newCmd, detach); err != nil {
		return errors.Wrapf(err, "failed run %s in container %s", newCmd, ctr.Digest)
	}
	return nil
}

// setNamespace calls setns syscall for set of flags. It changes
// current process namespace to namespace of another process which
// can be specified by pid.
//
// NOTE: A process may not be reassociated with a new mount namespace
// if it is multi-threaded. Changing the mount namespace requires that
// the caller possess both CAP_SYS_CHROOT and CAP_SYS_ADMIN capabilities
// in its own user namespace and CAP_SYS_ADMIN in the target mount namespace.
func setNamespace(pid int, flag int) error {
	nsBase := filepath.Join("/proc", strconv.Itoa(pid), "ns")
	ns := map[int]string{
		syscall.CLONE_NEWIPC: "ipc",
		syscall.CLONE_NEWNS:  "mnt",
		syscall.CLONE_NEWNET: "net",
		syscall.CLONE_NEWPID: "pid",
		syscall.CLONE_NEWUTS: "uts",
	}

	for k, v := range ns {
		if flag&k == 0 {
			continue
		}
		nsFile, err := os.Open(filepath.Join(nsBase, v))
		if err != nil {
			return errors.Wrapf(err, "can't open %s", nsFile)
		}

		if err := unix.Setns(int(nsFile.Fd()), k); err != nil {
			return errors.Wrapf(err, "can't setns to %s", v)
		}
	}

	return nil
}
