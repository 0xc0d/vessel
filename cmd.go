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

	flags := cmd.Flags()
	flags.StringP("name", "", "", "Container name")
	flags.StringP("host", "", "", "Container Hostname")
	flags.IntP("memory", "m", 100, "Limit memory access in MB")
	flags.IntP("swap", "s", 20, "Limit swap access in MB")
	flags.Float64P("cpus", "c", 2, "Limit CPUs")
	flags.IntP("pids", "p", 100, "Limit number of processes")
	flags.BoolP("detach", "d", false, "run command in the background")

	return cmd
}

func newForkCommand() *cobra.Command {
	ctr := new(container)
	cmd := &cobra.Command{
		Use:    "fork",
		//Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// Convert memory and swap limit to megabyte
			ctr.mem *= MB
			ctr.swap *= MB
			runFork(ctr, args[0], args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ctr.image, "image", "", "Container name")
	flags.StringVarP(&ctr.name, "name", "", "", "Container name")
	flags.StringVarP(&ctr.hostname, "host", "", "", "Container Hostname")
	flags.IntVarP(&ctr.mem, "memory", "m", 100, "Limit memory access in MB")
	flags.IntVarP(&ctr.swap, "swap", "s", 20, "Limit swap access in MB")
	flags.Float64VarP(&ctr.cpus, "cpus", "c", 2, "Limit CPUs")
	flags.IntVarP(&ctr.pids, "pids", "p", 128, "Limit number of processes")
	flags.BoolVarP(&ctr.detach, "detach", "d", false, "run command in the background")

	return cmd
}
