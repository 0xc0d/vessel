package cmd

import (
	"github.com/0xc0d/vessel/internal"
	"github.com/spf13/cobra"
)

// NewPsCommand implements and returns the ps command.
func NewPsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ps",
		Short:                 "List Containers",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		RunE:                  internal.Ps,
	}

	return cmd
}
