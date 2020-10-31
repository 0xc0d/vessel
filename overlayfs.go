package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func overlayFsMount(mountPoint string, imgLayers []string, readOnly bool) error {
	var upper, work []string
	if !readOnly {
		// Create Upper and Work Directories for writable mount
		parentDir := filepath.Dir(strings.TrimRight(mountPoint, "/"))
		upperDir := filepath.Join(parentDir, "diff")
		workDir := filepath.Join(parentDir, "word")
		if err := os.MkdirAll(upperDir, 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(workDir, 0755); err != nil {
			return err
		}
		upper = append(upper, upperDir)
		work = append(work, workDir)
	}

	opt := formatOverlayFsMountOption(imgLayers, upper, work)
	return syscall.Mount("none", mountPoint, "overlay", 0, opt)
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
