package cgroups

import (
	"bufio"
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
	Path      string
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

// NewCGroup creates an empty CGroups
func NewCGroup() *CGroups {
	return new(CGroups)
}

// SetPath sets Path for CGroups.
//
// NOTE: it does not require CGroup base path (/sys/fs/cgroup)
func (cg *CGroups) SetPath(path string) *CGroups {
	cg.Path = path
	return cg
}

// SetMemorySwapLimit sets memory and swap limit for CGroups.
//
// Zero or lower values consider as MAX.
func (cg *CGroups) SetMemorySwapLimit(memory, swap int) *CGroups {
	if memory > 1 {
		cg.mem = []byte(strconv.Itoa(memory))
		if swap > 1 {
			cg.memsw = []byte(strconv.Itoa(memory + swap))
		}
	}
	return cg
}

// SetCPULimit sets number of CPU for the CGroups.
func (cg *CGroups) SetCPULimit(quota float64) *CGroups {
	if int(quota) < runtime.NumCPU() && int(quota) > 0 {
		cg.cfsPeriod = []byte(strconv.Itoa(defaultCfsPeriod))
		cg.cfsQuota = []byte(strconv.Itoa(int(defaultCfsPeriod * quota)))
	}
	return cg
}

// SetProcessLimit sets maximum processes than can be created
// simultaneously in CGroups.
func (cg *CGroups) SetProcessLimit(number int) *CGroups {
	if number > 0 {
		cg.pids = []byte(strconv.Itoa(number))
	}
	return cg
}


// Load affects CGroups.
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

// Remove removes CGroups.
//
// It will only works if there is no process running in the CGroups
func (cg *CGroups) Remove() error {
	if cg.Path == "" {
		return errors.New("empty")
	}
	for _, c := range controllers {
		dir := filepath.Join(cgroupPath, c, cg.Path)
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

// GetPids returns slice of pids running on CGroups.
func (cg *CGroups) GetPids() ([]int, error) {
	var pids []int

	proc := filepath.Join(cgroupPath, controllers[0], cg.Path, procsFilename)
	procFile, err := os.Open(proc)
	if err != nil {
		return pids, err
	}
	defer procFile.Close()

	scanner := bufio.NewScanner(procFile)
	for scanner.Scan() {
		pid, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return pids, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (cg *CGroups) createControllersDir() error {
	for _, c := range controllers {
		dir := filepath.Join(cgroupPath, c, cg.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// enableReleaseAgent enables notify_on_release for CGroup.
func (cg *CGroups) enableReleaseAgent() error {
	for _, c := range controllers {
		file := filepath.Join(cgroupPath, c, cg.Path, releaseAgentFilename)
		if err := ioutil.WriteFile(file, []byte{'1'}, 0644); err != nil {
			return err
		}
	}
	return nil
}

// addProcess adds a pid into a CGroup.
func (cg *CGroups) addProcess(pid int) error {
	for _, c := range controllers {
		file := filepath.Join(cgroupPath, c, cg.Path, procsFilename)
		if err := ioutil.WriteFile(file, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CGroups) loadMemSwLimit() error {
	memoryLimitFile := filepath.Join(cgroupPath, "memory", cg.Path, memoryLimitFilename)
	memswLimitFile := filepath.Join(cgroupPath, "memory", cg.Path, memswLimitFilename)
	if err := ioutil.WriteFile(memoryLimitFile, cg.mem, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(memswLimitFile, cg.memsw, 0644)
}

func (cg *CGroups) loadCPULimit() error {
	cfsPeriodFile := filepath.Join(cgroupPath, "cpu", cg.Path, cpuPeriodFilename)
	cfsQuotaFile := filepath.Join(cgroupPath, "cpu", cg.Path, cpuQuotaFilename)
	if err := ioutil.WriteFile(cfsPeriodFile, cg.cfsPeriod, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(cfsQuotaFile, cg.cfsQuota, 0644)
}

func (cg *CGroups) loadProcessLimit() error {
	maxProcessFile := filepath.Join(cgroupPath, "pids", cg.Path, maxProcessFilename)
	return ioutil.WriteFile(maxProcessFile, cg.pids, 0644)
}
