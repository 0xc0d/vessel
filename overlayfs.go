package main

import (
	"os"
	"path/filepath"
	"strings"
)

func overlayFsMount(target string, imgLayers []string, readOnly bool) (func() error, error) {
	var upper, work []string
	if !readOnly {
		// Create Upper and Work Directories for writable mount
		parentDir := filepath.Dir(strings.TrimRight(target, "/"))
		upperDir := filepath.Join(parentDir, "diff")
		workDir := filepath.Join(parentDir, "word")
		if err := os.MkdirAll(upperDir, 0755); err != nil {
			message := "can't create overlay upper directory"
			return nil, errorWithMessage(err, message)
		}
		if err := os.MkdirAll(workDir, 0755); err != nil {
			message := "can't create overlay work directory"
			return nil, errorWithMessage(err, message)
		}
		upper = append(upper, upperDir)
		work = append(work, workDir)
	}

	opt := formatOverlayFsMountOption(imgLayers, upper, work)
	newMountPoint := mountPoint{
		source: "none",
		target: target,
		fsType: "overlay",
		flag:   0,
		option: opt,
	}

	return mount(newMountPoint)
}

func formatOverlayFsMountOption(lowerDir, upperDir, workDir []string) string {
	lower := "lowerdir="
	lower += strings.Join(lowerDir, ":")
	upper := "upperdir="
	upper += strings.Join(upperDir, ":")
	work := "workdir="
	work += strings.Join(workDir, ":")

	opt := strings.Join([]string{lower, upper, work}, ",")
	return opt
}
