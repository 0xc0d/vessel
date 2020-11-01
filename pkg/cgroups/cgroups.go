package cgroups

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	cgroupPath           = "/sys/fs/cgroup"
	releaseAgentFilename = "notify_on_release"
	procsFilename        = "cgroup.procs"
	memoryLimitFilename  = "memory.limit_in_bytes"
	memswLimitFilename   = "memory.memsw.limit_in_bytes"
	cpuQuotaFilename     = "cpu.cfs_quota_us"
	cpuPeriodFilename    = "cpu.cfs_period_us"
	maxProcessFilename   = "pids.max"

	defaultCfsPeriod = 100000
)

type CGroups struct {
	path      string
	mem       []byte
	memsw     []byte
	cfsPeriod []byte
	cfsQuota  []byte
	pids      []byte
}

var controllers = []string{
	"memory",
	"cpu",
	"pids",
}

func NewCGroup() *CGroups {
	return new(CGroups)
}

func (cg *CGroups) SetPath(path string) *CGroups {
	if path != "" {
		cg.path = path
	}
	return cg
}

func (cg *CGroups) SetMemorySwapLimit(memory, swap int) *CGroups {
	if memory > 1 {
		cg.mem = []byte(strconv.Itoa(memory))
		if swap > 1 {
			cg.memsw = []byte(strconv.Itoa(memory + swap))
		}
	}
	return cg
}

func (cg *CGroups) SetCPULimit(quota float64) *CGroups {
	if int(quota) < runtime.NumCPU() && int(quota) > 0 {
		cg.cfsPeriod = []byte(strconv.Itoa(defaultCfsPeriod))
		cg.cfsQuota = []byte(strconv.Itoa(int(defaultCfsPeriod * quota)))
	}
	return cg
}

func (cg *CGroups) SetProcessLimit(number int) *CGroups {
	if number > 0 {
		cg.pids = []byte(strconv.Itoa(number))
	}
	return cg
}

func (cg *CGroups) Load() error {
	if err := cg.createControllersDir(); err != nil {
		return err
	}
	if err := cg.enableReleaseAgent(); err != nil {
		return err
	}
	if err := cg.addProcess(os.Getpid()); err != nil {
		return err
	}
	if err := cg.loadMemSwLimit(); err != nil {
		return err
	}
	if err := cg.loadCPULimit(); err != nil {
		return err
	}
	if err := cg.loadProcessLimit(); err != nil {
		return err
	}

	return nil
}

func (cg *CGroups) Remove() error {
	if cg.path == "" {
		return errors.New("empty")
	}
	for _, c := range controllers {
		dir := filepath.Join(cgroupPath, c, cg.path)
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CGroups) createControllersDir() error {
	for _, c := range controllers {
		dir := filepath.Join(cgroupPath, c, cg.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CGroups) enableReleaseAgent() error {
	for _, c := range controllers {
		file := filepath.Join(cgroupPath, c, cg.path, releaseAgentFilename)
		if err := ioutil.WriteFile(file, []byte{'1'}, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CGroups) addProcess(pid int) error {
	for _, c := range controllers {
		file := filepath.Join(cgroupPath, c, cg.path, procsFilename)
		if err := ioutil.WriteFile(file, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CGroups) loadMemSwLimit() error {
	memoryLimitFile := filepath.Join(cgroupPath, "memory", cg.path, memoryLimitFilename)
	memswLimitFile := filepath.Join(cgroupPath, "memory", cg.path, memswLimitFilename)
	if err := ioutil.WriteFile(memoryLimitFile, cg.mem, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(memswLimitFile, cg.memsw, 0644)
}

func (cg *CGroups) loadCPULimit() error {
	cfsPeriodFile := filepath.Join(cgroupPath, "cpu", cg.path, cpuPeriodFilename)
	cfsQuotaFile := filepath.Join(cgroupPath, "cpu", cg.path, cpuQuotaFilename)
	if err := ioutil.WriteFile(cfsPeriodFile, cg.cfsPeriod, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(cfsQuotaFile, cg.cfsQuota, 0644)
}

func (cg *CGroups) loadProcessLimit() error {
	maxProcessFile := filepath.Join(cgroupPath, "pids", cg.path, maxProcessFilename)
	return ioutil.WriteFile(maxProcessFile, cg.pids, 0644)
}
