package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
		SilenceUsage:          true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if os.Getuid() != 0 { // not root
				message := fmt.Sprintf("%s should be run as root", args[0])
				return errorWithMessage(ErrNotPermitted, message)
			}
			return nil
		},
		RunE: runRun,
	}

	flags := cmd.Flags()
	flags.StringP("name", "", "", "Container name")
	flags.StringP("host", "", "", "Container Hostname")
	flags.IntP("memory", "m", 100, "Limit memory access in MB")
	flags.IntP("swap", "s", 20, "Limit swap access in MB")
	flags.Float64P("cpus", "c", 2, "Limit CPUs")
	flags.IntP("pids", "p", 128, "Limit number of processes")
	flags.BoolP("detach", "d", false, "run command in the background")
	return cmd
}

func newForkCommand() *cobra.Command {
	ctr := new(container)
	cmd := &cobra.Command{
		Use:           "fork",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Convert memory and swap limit to megabyte
			ctr.mem *= MB
			ctr.swap *= MB
			return runFork(ctr, args)
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
