package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// runRun runs a command inside a new container
func runRun(cmd *cobra.Command, args []string) {
	img, err := getImage(args[0])
	CheckErr(err)

	ctr := new(container)
	flags := cmd.Flags()
	flags.StringVarP(&ctr.name, "name", "", "", "Container name")
	flags.StringVarP(&ctr.hostname, "host", "", "", "Container Hostname")
	flags.IntVarP(&ctr.mem, "memory", "m", 100, "Limit memory access in MB")
	flags.IntVarP(&ctr.swap, "swap", "s", 20, "Limit swap access in MB")
	flags.Float64VarP(&ctr.cpus, "cpus", "c", 2, "Limit CPUs")
	flags.IntVarP(&ctr.pids, "pids", "p", 128, "Limit number of processes")
	flags.BoolVarP(&ctr.detach, "detach", "d", false, "run command in the background")

	ctr.setDigest()
	Must(ctr.mountFromImage(img))
	defer ctr.unmountFs()

	options := []string{fmt.Sprintf("--container=%s", ctr.digest)}
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Value.String() != "" {
			options = append(options, fmt.Sprintf("--%s=%v", flag.Name, flag.Value))
		}
	})

	// Add environment values to flags
	imgConfig, err := img.ConfigFile()
	CheckErr(err)
	env := strings.Join(imgConfig.Config.Env, ",")
	options = append(options, fmt.Sprintf("--environments=%q", env))

	commandToExec := args[1:]
	newArgs := append([]string{"fork"}, options...)
	newArgs = append(newArgs, commandToExec...)

	newCmd := exec.Command("/proc/self/exe", newArgs...)
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC,
	}

	newCmd.Run()
	newCmd.Wait()
}

func runFork(ctr *container, arg []string) {
	rootDir := filepath.Join(CtrDir, ctr.digest, "mnt")
	if err := syscall.Chroot(rootDir); err != nil {
		log.Println(err, rootDir)
	}
	Must(os.Chdir("/"))
	//fmt.Println(os.MkdirAll("/sys", 0755))
	//fmt.Println(os.MkdirAll("/proc", 0755))
	Must(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("proc", 0)
	//Must(syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""))
	//Must(syscall.Mount("tmpfs", "dev", "tmpfs", 0, ""))
	fmt.Println(syscall.Mount("sysfs", "sys", "sysfs", 0, ""))
	defer syscall.Unmount("sys", 0)

	ctr.setHostname()

	newCmd := exec.Command("/bin/sh")
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.Env = ctr.env
	newCmd.Run()
}
