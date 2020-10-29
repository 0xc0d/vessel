package main

const (
	VERSION = "0.1"
)

func main() {
	a := newVesselCommand()
	a.UsageFunc()
	a.AddCommand(newRunCommand())
	must(a.Execute())
}