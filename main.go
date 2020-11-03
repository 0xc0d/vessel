package main

import (
	"fmt"
	"github.com/0xc0d/vessel/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewVesselCommand()
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewForkCommand())
	rootCmd.AddCommand(cmd.NewExecCommand())
	rootCmd.AddCommand(cmd.NewPsCommand())
	rootCmd.AddCommand(cmd.NewImagesCommand())
	err := rootCmd.Execute()
	fmt.Fprintln(os.Stderr, err)
}