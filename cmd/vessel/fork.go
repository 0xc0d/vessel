package vessel

import (
	"errors"
	"github.com/0xc0d/vessel/internal"
	"github.com/0xc0d/vessel/pkg/container"
	"github.com/spf13/cobra"
)

var ErrEmptyRootFS = errors.New("root flag is required")

func NewForkCommand() *cobra.Command {
	ctr := container.NewContainer()
	var detach bool
	cmd := &cobra.Command{
		Use:    "fork",
		//Hidden: true,
		//SilenceUsage:  true,
		//SilenceErrors: true,
		PreRunE: isRoot,
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
	flags.StringVar(&ctr.Name, "name", "", "")
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
