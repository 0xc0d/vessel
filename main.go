package main

const (
	VERSION = "0.1"
)

func main() {
	a := newVesselCommand()
	a.AddCommand(newRunCommand())
	a.Execute()
}
