package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

const (
	cgroupPath           = "/sys/fs/cgroup"
	vessel               = "vessel"
	releaseAgentFilename = "notify_on_release"
	procsFilename        = "cgroup.procs"
	memoryLimitFilename  = "memory.limit_in_bytes"
	memswLimitFilename   = "memory.memsw.limit_in_bytes"
	cpuQuotaFilename     = "cpu.cfs_quota_us"
	cpuPeriodFilename    = "cpu.cfs_period_us"
	maxProcessFilename   = "pid.max"

	defaultCfsPeriod   = 100000
	defaultMaxProcess  = 1 << 7
	defaultMemoryLimit = 1 << 7 * MB
	defaultSwapLimit   = 1 << 5 * MB
)

type cgroup struct {
	ctrid     string
	mem       []byte
	memsw     []byte
	cfsPeriod []byte
	cfsQuota  []byte
	pids      []byte
}

var controller = []string{
	"memory",
	"cpu",
	"pids",
}

func newCGroup(ctr *container) *cgroup {
	if int(ctr.cpus) > runtime.NumCPU() {
		ctr.cpus = -1 // no limit
	}
	if ctr.pids < 1 {
		ctr.pids = defaultMaxProcess
	}
	if ctr.mem < 1 {
		ctr.mem = defaultMemoryLimit
	}
	if ctr.swap < 1 {
		ctr.swap = defaultSwapLimit
	}

	memsw := ctr.mem + ctr.swap
	cg := &cgroup{
		ctrid:     ctr.digest,
		mem:       []byte(strconv.Itoa(ctr.mem)),
		memsw:     []byte(strconv.Itoa(memsw)),
		cfsPeriod: []byte(strconv.Itoa(defaultCfsPeriod)),
		cfsQuota:  []byte(strconv.Itoa(int(defaultCfsPeriod * ctr.cpus))),
		pids:      []byte(strconv.Itoa(ctr.pids)),
	}

	return cg
}

func (cg *cgroup) Load() error {
	if err := cg.createControllersDir(); err != nil {
		return err
	}
	if err := cg.enableReleaseAgent(); err != nil {
		return err
	}
	if err := cg.moveProcess(os.Getpid()); err != nil {
		return err
	}
	if err := cg.setMemSwLimit(); err != nil {
		return err
	}
	if err := cg.setCPULimit(); err != nil {
		return err
	}
	if err := cg.setProcessLimit(); err != nil {
		return err
	}

	return nil
}

func (cg *cgroup) Remove() error {
	for _, c := range controller {
		dir := filepath.Join(cgroupPath, c, vessel, cg.ctrid)
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func (cg *cgroup) createControllersDir() error {
	for _, c := range controller {
		dir := filepath.Join(cgroupPath, c, vessel, cg.ctrid)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (cg *cgroup) enableReleaseAgent() error {
	for _, c := range controller {
		file := filepath.Join(cgroupPath, c, vessel, cg.ctrid, releaseAgentFilename)
		if err := ioutil.WriteFile(file, []byte{'1'}, 0700); err != nil {
			return err
		}
	}
	return nil
}

func (cg *cgroup) moveProcess(pid int) error {
	for _, c := range controller {
		file := filepath.Join(cgroupPath, c, vessel, cg.ctrid, procsFilename)
		if err := ioutil.WriteFile(file, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (cg *cgroup) setMemSwLimit() error {
	memoryLimitFile := filepath.Join(cgroupPath, "memory", vessel, cg.ctrid, memoryLimitFilename)
	memswLimitFile := filepath.Join(cgroupPath, "memory", vessel, cg.ctrid, memswLimitFilename)
	if err := ioutil.WriteFile(memoryLimitFile, cg.mem, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(memswLimitFile, cg.memsw, 0644)
}

func (cg *cgroup) setCPULimit() error {
	cfsPeriodFile := filepath.Join(cgroupPath, "cpu", vessel, cg.ctrid, cpuPeriodFilename)
	cfsQuotaFile := filepath.Join(cgroupPath, "cpu", vessel, cg.ctrid, cpuQuotaFilename)
	if err := ioutil.WriteFile(cfsPeriodFile, cg.cfsPeriod, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(cfsQuotaFile, cg.cfsQuota, 0644)
}

func (cg *cgroup) setProcessLimit() error {
	maxProcessFile := filepath.Join(cgroupPath, "pids", vessel, cg.ctrid, maxProcessFilename)
	return ioutil.WriteFile(maxProcessFile, cg.pids, 0644)
}
