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
	ops := new(containerOptions)

	cmd := &cobra.Command{
		Use:                   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short:                 "Run a command inside a new container.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			src := args[0]
			ops.image, err = newImage(src)
			if err != nil {
				cmd.SilenceUsage = true
				return ErrRepoNotExist
			}

			// Convert memory and swap limit to megabyte
			ops.mem *= MB
			ops.swap *= MB
			return
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				Cmd string
				Args []string
			)
			if len(args) > 1 {
				Cmd = args[1]
			}
			if len(args) > 2 {
				Args = args[2:]
			}

			runRun(ops, Cmd, Args...)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&ops.name, "name", "", "", "Container name")
	flags.StringVarP(&ops.hostname, "host", "", "", "Container Hostname")
	flags.IntVarP(&ops.mem, "memory", "m", 100, "Limit memory access in MB")
	flags.IntVarP(&ops.swap, "swap", "s", 20, "Limit swap access in MB")
	flags.Float64VarP(&ops.cpus, "cpus", "c", 2, "Limit CPUs")
	flags.IntVarP(&ops.pid, "pids", "p", 100, "Limit number of processes")
	flags.BoolVarP(&ops.detach, "detach", "d", false, "run command in the background")

	return cmd
}