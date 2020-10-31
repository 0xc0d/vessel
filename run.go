package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// runRun runs a command inside a new container
func runRun(cmd *cobra.Command, args []string) error {
	img, err := getImage(args[0])
	if err != nil {
		message := fmt.Sprintf("local/remote image %s not found", args[0])
		return errorWithMessage(err, message)
	}

	// create container from flags and set a random hash for it
	ctr := new(container)
	ctr.setDigest()

	// mount image layer for container use
	unmounter, err := ctr.mountFromImage(img)
	if err != nil {
		return errorWithMessage(err, "can't mount image filesystem")
	}
	defer unmounter()

	// Format fork options
	var options []string
	{
		options = append(options, fmt.Sprintf("--container=%s", ctr.digest))
		flags := cmd.Flags()
		flags.VisitAll(func(flag *pflag.Flag) {
			if flag.Value.String() != "" {
				options = append(options, fmt.Sprintf("--%s=%v", flag.Name, flag.Value))
			}
		})
		// Add environment values to flags
		imgConfig, err := img.ConfigFile()
		if err != nil {
			return errorWithMessage(err, "can't get image config")
		}
		env := strings.Join(imgConfig.Config.Env, ",")
		options = append(options, fmt.Sprintf("--environments=%q", env))
	}

	commandToExec := args[1:]
	newArgs := append([]string{"fork"}, options...)
	newArgs = append(newArgs, commandToExec...)

	newCmd := exec.Command("/proc/self/exe", newArgs...)
	newCmd.Stdin, newCmd.Stdout, newCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	var flag uintptr
	flag = syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC
	newCmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: flag}

	if err := newCmd.Run(); err != nil {
		return errorWithMessage(err, "failed run fork process")
	}

	// remove CGroups created by fork process
	if err := ctr.removeCGroups(); err != nil {
		return errorWithMessage(err, "can't remove cgroup files")
	}

	return nil
}

func runFork(ctr *container, arg []string) error {
	ctr.setHostname()

	if err := ctr.loadCGroups(); err != nil {
		return errorWithMessage(err, "can't initialize cgroups")
	}

	rootDir := filepath.Join(CtrDir, ctr.digest, "mnt")
	if err := chroot(rootDir); err != nil {
		return err
	}

	mountPoints := []mountPoint{
		{source: "proc", target: "proc", fsType: "proc", flag: 0, option: ""},
		{source: "sysfs", target: "sys", fsType: "sysfs", flag: 0, option: ""},
	}
	unmounter, err := mount(mountPoints...)
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
	newCmd.Env = ctr.env
	return newCmd.Run()
}

func chroot(root string) error {
	if err := syscall.Chroot(root); err != nil {
		message := fmt.Sprintf("can't change root to %s", root)
		return errorWithMessage(err, message)
	}
	return os.Chdir("/")
}
