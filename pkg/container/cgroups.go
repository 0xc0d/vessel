package container

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