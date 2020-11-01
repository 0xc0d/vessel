package vessel

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var ErrNotPermitted = errors.New("operation not permitted")

// NewVesselCommand returns the root cobra.Command for Vessel.
func NewVesselCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "vessel [OPTIONS] COMMAND",
		Short:                 "A tiny tool for managing containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

func isRoot(_ *cobra.Command, _ []string) error {
	if os.Getuid() != 0 {
		return ErrNotPermitted
	}
	return nil
}
