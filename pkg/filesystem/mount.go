package filesystem

import (
	"github.com/pkg/errors"
	"os"
	"syscall"
)

type MountPoint struct {
	Source string
	Target string
	Type   string
	Flag   uintptr
	Option string
}

type Unmounter func() error

func Mount(mountPoints...MountPoint) (Unmounter, error) {
	unmounter := func() error {
		for _, p := range mountPoints {
			if err := syscall.Unmount(p.Target, 0); err != nil {
				return errors.Wrapf(err, "unable to umount %q", p.Target)
			}
		}
		return nil
	}

	for _, p := range mountPoints {
		if err := os.MkdirAll(p.Target, 0755); err != nil {
			return unmounter, errors.Wrapf(err, "can't create %s directory", p.Target)
		}
		if err := syscall.Mount(p.Source, p.Target, p.Type, p.Flag, p.Option); err != nil {
			return unmounter, errors.Wrapf(err, "unable to mount %s to %s", p.Source, p.Target)
		}
	}

	return unmounter, nil
}