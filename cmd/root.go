package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

const (
	layersPath     = "/var/lib/vessel/images/layers"
	containersPath = "/var/run/vessel/containers"
	netnsPath      = "/var/run/vessel/netns"
)

var ErrNotPermitted = errors.New("operation not permitted")

// Make vessel directories first.
func init() {
	os.MkdirAll(netnsPath, 0700)
	os.MkdirAll(layersPath, 0700)
	os.MkdirAll(containersPath, 0700)
}

// NewVesselCommand returns the root cobra.Command for Vessel.
func NewVesselCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "vessel [OPTIONS] COMMAND",
		Short:                 "A tiny tool for managing containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE:     isRoot,
	}

	return cmd
}

// isRoot implements a cobra acceptable function and
// returns ErrNotPermitted if user is not root.
func isRoot(_ *cobra.Command, _ []string) error {
	if os.Getuid() != 0 {
		return ErrNotPermitted
	}
	return nil
}
