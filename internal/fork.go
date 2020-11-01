package internal

import (
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"syscall"
)

func Fork(ctr *container.Container, arg []string, detach bool) error {
	ctr.SetHostname()

	if err := ctr.LoadCGroups(); err != nil {
		return errors.Wrap(err, "can't initialize cgroups")
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

	newCmd := exec.Command(arg[0])
	if len(arg) > 1 {
		newCmd.Args = arg[1:]
	}
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.Env = ctr.Config.Env
	return runCommand(newCmd, detach)
}

func changeRoot(root, workdir string) error {
	if err := syscall.Chroot(root); err != nil {
		return errors.Wrapf(err, "can't change root to %s", root)
	}
	if workdir != "" {
		return os.Chdir(workdir)
	}
	return os.Chdir("/")
}

func runCommand(cmd *exec.Cmd, detach bool) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if !detach {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}
