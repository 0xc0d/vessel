package main

import (
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
	cmd := &cobra.Command{
		Use:                   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short:                 "Run a command inside a new container.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return
		},
		Run: runRun,
	}

	return cmd
}

func newForkCommand() *cobra.Command {
	ctr := new(container)
	cmd := &cobra.Command{
		Use:    "fork",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// Convert memory and swap limit to megabyte
			ctr.mem *= MB
			ctr.swap *= MB
			runFork(ctr, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ctr.digest, "container", "", "")
	flags.StringSliceVar(&ctr.env, "environments", []string{}, "")
	flags.StringVar(&ctr.name, "name", "", "")
	flags.StringVar(&ctr.hostname, "host", "", "")
	flags.IntVar(&ctr.mem, "memory", 100, "")
	flags.IntVar(&ctr.swap, "swap", 20, "")
	flags.Float64Var(&ctr.cpus, "cpus", 2, "")
	flags.IntVar(&ctr.pids, "pids", 128, "")
	flags.BoolVar(&ctr.detach, "detach", false, "")

	return cmd
}
