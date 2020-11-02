package internal

import (
	"fmt"
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Run runs a command inside a new container
func Ps(_ *cobra.Command, _ []string) error {
	allCtr, err := container.GetAllContainer()
	if err != nil {
		return err
	}

	pPrint(allCtr)
	return nil
}

func pPrint(ctrs []*container.Container) {
	fmt.Println("CONTAINER ID\t\tIMAGE       \t\tCOMMAND")
	for _, ctr := range ctrs {
		pids, err := ctr.GetPids()
		if err != nil {
			continue
		}

		// pid[0] is the reexec commmand
		cmd, err := getCmdlineById(pids[1])
		if err != nil {
			continue
		}
		image := strings.TrimLeft(ctr.Config.Image, "sha256:")
		fmt.Printf("%.12s\t\t%.12s\t\t%s\n", ctr.Digest, image, cmd)
	}
}

func getCmdlineById(pid int) (string, error) {
	cmdline, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	return string(cmdline), err
}
