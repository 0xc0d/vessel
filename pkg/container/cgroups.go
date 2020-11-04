package container

import (
	"github.com/0xc0d/vessel/pkg/cgroups"
	"path/filepath"
)

// LoadCGroups loads CGroups for container.
func (c *Container) LoadCGroups() error {
	cg := cgroups.NewCGroup()
	cg.SetPath(filepath.Join("vessel", c.Digest)).
		SetMemorySwapLimit(c.mem, c.swap).
		SetCPULimit(c.cpus).
		SetProcessLimit(c.pids)

	err := cg.Load()
	if err != nil {
		return err
	}
	pids, err := cg.GetPids()
	c.Pids = pids
	return err
}

// RemoveCGroups removes CGroups file for container.
// It only function if the container is not running.
func (c *Container) removeCGroups() error {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", c.Digest),
	}
	return cg.Remove()
}

// SetMemorySwapLimit sets Container's memory and swap limitation in MegaByte.
func (c *Container) SetMemorySwapLimit(memory, swap int) *Container {
	c.mem = memory * MB
	c.swap = swap * MB
	return c
}

// SetCPULimit sets Container number of CPUs.
func (c *Container) SetCPULimit(cpus float64) *Container {
	c.cpus = cpus
	return c
}

// SetProcessLimit sets maximum simultaneous process for Container.
func (c *Container) SetProcessLimit(pids int) *Container {
	c.pids = pids
	return c
}

// GetPids returns slice of pid running inside Container.
//
// NOTE: First element [0], is the fork process.
func (c *Container) GetPids() ([]int, error) {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", c.Digest),
	}
	pids, err := cg.GetPids()
	return pids, err
}

// getPidsByDigest returns slice of pid running inside a Container.
// Container should be specified by its digest.
func getPidsByDigest(digest string) ([]int, error) {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", digest),
	}
	pids, err := cg.GetPids()
	return pids, err
}
