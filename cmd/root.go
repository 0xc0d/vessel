package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

const (
	vesselDir = "/var/lib/vessel"
	ImgDir    = "/var/lib/vessel/images"
	LyrDir    = "/var/lib/vessel/images/layers"
	CtrDir    = "/var/lib/vessel/containers"
)

var ErrNotPermitted = errors.New("operation not permitted")

func init() {
	os.MkdirAll(vesselDir, 0711)
	os.MkdirAll(ImgDir, 0700)
	os.MkdirAll(LyrDir, 0700)
	os.MkdirAll(CtrDir, 0700)
}

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
