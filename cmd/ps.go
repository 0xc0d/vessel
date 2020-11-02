package cmd

import (
	"github.com/0xc0d/vessel/internal"
	"github.com/spf13/cobra"
)

// NewRunCommand implements and returns the run command.
func NewPsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ps",
		Short:                 "List Containers",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args: cobra.MaximumNArgs(0),
		PreRunE:               isRoot,
		RunE:                  internal.Ps,
	}

	return cmd
}
