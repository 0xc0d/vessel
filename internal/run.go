package internal

import (
	"fmt"
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/0xc0d/vessel/pkg/image"
	"github.com/0xc0d/vessel/pkg/reexec"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
	"syscall"
)

// Run runs a command inside a new container
func Run(cmd *cobra.Command, args []string) error {
	ctr := container.NewContainer()

	img, err := getImage(args[0])
	if err != nil {
		return err
	}
	// Mount image layer for container use
	unmount, err := ctr.MountFromImage(img)
	if err != nil {
		return errors.Wrap(err, "can't Mount image filesystem")
	}

	// defer container cleanup function
	defer func() {
		unmount()
		ctr.Remove()
		ctr.RemoveCGroups()
	}()

	// Format fork options
	options := append([]string{}, rawFlags(cmd.Flags())...)
	options = append(options, fmt.Sprintf("--root=%s", ctr.RootFS))
	options = append(options, fmt.Sprintf("--container=%s", ctr.Digest))

	newArgs := []string{"fork"}
	newArgs = append(newArgs, options...)
	newArgs = append(newArgs, args[1:]...)

	newCmd := reexec.Command(newArgs...)
	newCmd.Stdin, newCmd.Stdout, newCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	newCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	if err := newCmd.Run(); err != nil {
		return errors.Wrap(err, "failed run fork process")
	}

	return nil
}

func rawFlags(flags *pflag.FlagSet) []string {
	var flagList []string
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Value.String() == "" {
			return
		}
		flagList = append(flagList, fmt.Sprintf("--%s=%v", flag.Name, flag.Value))
	})
	return flagList
}

func getImage(name string) (v1.Image, error) {
	nameWithTag := withTag(name)
	img, err := image.NewImage(nameWithTag)
	if err != nil {
		return img, errors.Wrapf(err, "Can't pull %q", nameWithTag)
	}
	if image.Exists(img) {
		return img, nil
	}
	fmt.Printf("Unable to find image %q locally\n", nameWithTag)
	return img, image.Download(img)
}

func withTag(name string) string {
	nameTag := strings.Split(name, ":")
	if len(nameTag) != 2 || nameTag[1] == "" {
		nameTag = append(nameTag, "latest")
	}
	return strings.Join(nameTag, ":")
}