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
	c.Pid, err = cg.GetPid()
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

func GetPidByDigest(digest string) (int, error) {
	cg := &cgroups.CGroups{
		Path: filepath.Join("vessel", digest),
	}
	return cg.GetPid()
}