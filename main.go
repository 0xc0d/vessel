package main

import (
	"github.com/0xc0d/vessel/cmd"
)

func main() {
	rootCmd := cmd.NewVesselCommand()
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewForkCommand())
	rootCmd.Execute()
}