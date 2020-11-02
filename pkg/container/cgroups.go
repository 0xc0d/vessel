package container

import (
	"github.com/0xc0d/vessel/pkg/cgroups"
	"path/filepath"
)

func (c *Container) LoadCGroups() error {
	cg := cgroups.NewCGroup()
	cg.SetPath(filepath.Join("vessel", c.Digest)).
		SetMemorySwapLimit(c.mem*MB, c.swap*MB).
		SetCPULimit(c.cpus).
		SetProcessLimit(c.pids)

	err := cg.Load()
	if err != nil {
		return err
	}
	pids, err := cg.GetPids()
	c.Pid = pids[0]
	return err
}

func (c *Container) RemoveCGroups() error {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", c.Digest),
	}
	return cg.Remove()
}

func (c *Container) SetMemorySwapLimit(memory, swap int) *Container {
	c.mem = memory
	c.swap = swap
	return c
}

func (c *Container) SetCPULimit(cpus float64) *Container {
	c.cpus = cpus
	return c
}

func (c *Container) SetProcessLimit(pids int) *Container {
	c.pids = pids
	return c
}

func (c *Container) GetPids() ([]int, error) {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", c.Digest),
	}
	pids, err := cg.GetPids()
	return pids, err
}

func GetPidsByDigest(digest string) ([]int, error) {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", digest),
	}
	pids, err := cg.GetPids()
	return pids, err
}