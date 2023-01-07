package cmd

import (
	"github.com/0xc0d/vessel/internal"
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/spf13/cobra"
)

// NewForkCommand implements and returns fork command.
// fork command is called by reexec to apply namespaces.
//
// It is a hidden command and requires root path and
// container id to run.
func NewForkCommand() *cobra.Command {
	ctr := container.NewContainer()
	var detach bool
	cmd := &cobra.Command{
		Use:          "fork",
		Hidden:       true,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctr.LoadConfig(); err != nil {
				return err
			}
			return internal.Fork(ctr, args, detach)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ctr.Digest, "container", "", "")
	flags.StringVar(&ctr.RootFS, "root", "", "")
	flags.StringVar(&ctr.Config.Hostname, "host", "", "")
	flags.BoolVar(&detach, "detach", false, "")
	mem := flags.Int("memory", 100, "")
	swap := flags.Int("swap", 20, "")
	cpu := flags.Float64("cpus", 1, "")
	pids := flags.Int("pids", 128, "")
	ctr.SetMemorySwapLimit(*mem, *swap)
	ctr.SetCPULimit(*cpu)
	ctr.SetProcessLimit(*pids)

	cmd.MarkFlagRequired("root")
	cmd.MarkFlagRequired("container")
	return cmd
}
