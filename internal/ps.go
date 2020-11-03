package internal

import (
	"bytes"
	"fmt"
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Ps gets all running containers and prints them.
func Ps(_ *cobra.Command, _ []string) error {
	allCtr, err := container.GetAllContainers()
	if err != nil {
		return err
	}

	pPrintContainer(allCtr)
	return nil
}

func pPrintContainer(ctrs []*container.Container) {
	fmt.Println("CONTAINER ID\t\tIMAGE       \t\tCOMMAND")
	for _, ctr := range ctrs {
		pids, err := ctr.GetPids()
		if err != nil || len(pids) < 2 {
			continue
		}

		// pid[0] is the reexec command
		// pid[1] is the init command
		cmd, err := cmdlineById(pids[1])
		if err != nil {
			continue
		}
		img := strings.TrimLeft(ctr.Config.Image, "sha256:")
		fmt.Printf("%.12s\t\t%.12s\t\t%.40q\n", ctr.Digest, img, cmd)
	}
}

// cmdlineById returns command associated with a pid.
func cmdlineById(pid int) (string, error) {
	cmdline, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	cmdline = bytes.ReplaceAll(cmdline, []byte{0}, []byte{' '})
	cmdline = bytes.TrimSpace(cmdline)
	return string(cmdline), err
}
