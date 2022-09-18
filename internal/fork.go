package internal

import (
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"syscall"
)

// Fork will call by Run. It is a hack to fork a whole new Go process
// inside a new namespace.
//
// If detach was enable function returns immediately after starting
// the command and never wait for result
func Fork(ctr *container.Container, args []string, detach bool) error {
	ctr.SetHostname()
	// set network
	unset, err := ctr.SetNetworkNamespace()
	if err != nil {
		return errors.Wrap(err, "can't set network namespace")
	}
	defer unset()

	if err := ctr.LoadCGroups(); err != nil {
		return errors.Wrap(err, "can't initialize cgroups")
	}
	if err := changeRoot(ctr.RootFS, ctr.Config.WorkingDir); err != nil {
		return err
	}

	// Mount necessaries
	mountPoints := []filesystem.MountOption{
		{Source: "proc", Target: "proc", Type: "proc"},
		{Source: "sysfs", Target: "sys", Type: "sysfs"},
	}
	unmount, err := filesystem.Mount(mountPoints...)
	if err != nil {
		return err
	}
	defer unmount()

	command, argv := cmdAndArgs(ctr.Config.Cmd)
	if len(args) > 0 {
		command, argv = cmdAndArgs(args)
	}
	newCmd := exec.Command(command, argv...)
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.Env = ctr.Config.Env
	return runCommand(newCmd, detach)
}

// changeRoot calls chroot syscall for the given root filesystem and will
// change working directory into workdir
func changeRoot(root, workdir string) error {
	if err := syscall.Chroot(root); err != nil {
		return errors.Wrapf(err, "can't change root to %s", root)
	}
	if workdir == "" {
		workdir = "/"
	}
	return os.Chdir(workdir)
}

// runCommand runs a command and wait or not wait for it based on detach.
func runCommand(cmd *exec.Cmd, detach bool) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if detach {
		return cmd.Process.Release()
	}
	
	return cmd.Wait()
}

// cmdAndArgs separate command (args[0]) and its argv.
func cmdAndArgs(args []string) (command string, argv []string) {
	if len(args) == 0 {
		return
	}
	command = args[0]
	argv = args[1:]
	return
}
