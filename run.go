package main

import "fmt"

// runRun runs a command inside a new container
func runRun(opts *containerOptions, command string, args... string) {
	if !imageExists(opts.image) {
		Must(imageDownload(opts.image))
	}
	fmt.Println("running", command, args)
}