package main

import "github.com/0xc0d/vessel/cmd/vessel"

func main() {
	rootCmd := vessel.NewVesselCommand()
	rootCmd.AddCommand(vessel.NewRunCommand())
	rootCmd.AddCommand(vessel.NewForkCommand())
	rootCmd.Execute()
}