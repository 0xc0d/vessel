package filesystem

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

// OverlayMount mounts a list of source directories to a target.
//
// If readOnly was true the target will be read-only.
func OverlayMount(target string, src []string, readOnly bool) (Unmounter, error) {
	var upper, work []string
	if !readOnly {
		// Create Upper and Work Directories for writable Mount
		parentDir := filepath.Dir(strings.TrimRight(target, "/"))
		upperDir := filepath.Join(parentDir, "diff")
		workDir := filepath.Join(parentDir, "work")
		if err := os.MkdirAll(upperDir, 0700); err != nil {
			return nil, errors.Wrap(err, "can't create overlay upper directory")
		}
		if err := os.MkdirAll(workDir, 0700); err != nil {
			return nil, errors.Wrap(err, "can't create overlay work directory")
		}
		upper = append(upper, upperDir)
		work = append(work, workDir)
	}

	opt := formatOverlayFsMountOption(src, upper, work)
	newMountOption := MountOption{
		Source: "none",
		Target: target,
		Type:   "overlay",
		Flag:   0,
		Option: opt,
	}

	return Mount(newMountOption)
}

// formatOverlayFsMountOption returns formatted overlayFS mount option.
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
