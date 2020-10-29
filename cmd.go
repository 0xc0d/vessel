package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// newVesselCommand returns the root cobra.Command for Vessel.
func newVesselCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "vessel [OPTIONS] COMMAND",
		Short:                 "A tiny tool for managing containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		Version:               VERSION,
	}

	return cmd
}

// newRunCommand implements and returns the run command.
func newRunCommand() *cobra.Command {
	ops := new(containerOptions)

	cmd := &cobra.Command{
		Use:                   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short:                 "Run a command inside a new container.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "ubuntu" {
				cmd.SilenceUsage = true
				return fmt.Errorf("can't find image %s", args[0])
			}
			ops.mem *= MB
			ops.swap *= MB
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(ops)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&ops.name, "name", "n", "", "Container name")
	flags.StringVarP(&ops.hostname, "host", "h", "", "Container Hostname")
	flags.IntVarP(&ops.mem, "memory", "m", 100, "Limit memory access in MB")
	flags.IntVarP(&ops.swap, "swap", "s", 20, "Limit swap access in MB")
	flags.Float64VarP(&ops.cpus, "cpus", "c", 2, "Limit CPUs")
	flags.IntVarP(&ops.pid, "pids", "p", 100, "Limit number of processes")
	flags.BoolVarP(&ops.detach, "detach", "d", false, "run command in the background")

	return cmd
}