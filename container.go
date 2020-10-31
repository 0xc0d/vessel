package main

import (
	"crypto/sha256"
	"fmt"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type container struct {
	digest   string
	name     string
	hostname string
	env      []string
	mem      int
	swap     int
	pids     int
	cpus     float64
	detach   bool
	once     sync.Once
}

const (
	CtrDir = "/home/aggy/var/lib/vessel/containers"
)

func (c *container) setDigest() {
	c.once.Do(func() {
		if len(c.digest) > 0 {
			return
		}
		rand.Seed(time.Now().Unix())
		randBuffer := make([]byte, 32)
		rand.Read(randBuffer)
		sha := sha256.New().Sum(randBuffer)
		c.digest = fmt.Sprintf("%x", sha)[:64]
	})
}

func (c *container) setHostname() {
	if c.hostname == "" {
		c.hostname = c.digest[:12]
	}
	Must(syscall.Sethostname([]byte(c.hostname)))
}

func (c *container) loadCGroups() error {
	cg := newCGroup()
	cg.setPath(filepath.Join("vessel", c.digest))
	cg.setMemorySwapLimit(c.mem, c.swap)
	cg.setCPULimit(c.cpus)
	cg.setProcessLimit(c.pids)
	return cg.Load()
}

func (c *container) removeCGroups() error {
	cg := newCGroup()
	cg.setPath(filepath.Join("vessel", c.digest))
	return cg.Remove()
}

func (c *container) mountFromImage(img v1.Image) error {
	mountPoint := filepath.Join(CtrDir, c.digest, "mnt")
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return err
	}
	layers, err := getImageLayers(img)
	if err != nil {
		return err
	}
	if err := overlayFsMount(mountPoint, layers, false); err != nil {
		return err
	}
	return nil
}

func (c *container) unmountFs() error {
	containerPath := filepath.Join(CtrDir, c.digest)
	mountPoint := filepath.Join(containerPath, "mnt")
	if err := syscall.Unmount(mountPoint, 0); err != nil {
		return err
	}
	return os.RemoveAll(containerPath)
}
