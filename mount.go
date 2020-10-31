package main

import (
	"fmt"
	"os"
	"syscall"
)

type mountPoint struct {
	source string
	target string
	fsType string
	flag   uintptr
	option string
}

func mount(mountPoints... mountPoint) (func() error, error) {
	unmounter := func() error {
		for _, p := range mountPoints {
			if err := syscall.Unmount(p.target, 0); err != nil {
				message := fmt.Sprintf("unable to umount %s", p.target)
				return errorWithMessage(err, message)
			}
		}
		return nil
	}

	for _, p := range mountPoints {
		if err := os.MkdirAll(p.target, 0755); err != nil {
			message := fmt.Sprintf("can't create %s directory", p.target)
			return unmounter, errorWithMessage(err, message)
		}
		if err := syscall.Mount(p.source, p.target, p.fsType, p.flag, p.option); err != nil {
			message := fmt.Sprintf("unable to mount %s to %s", p.source, p.target)
			return unmounter, errorWithMessage(err, message)
		}
	}

	return unmounter, nil
}