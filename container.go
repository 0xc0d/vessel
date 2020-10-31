package main

import (
	"crypto/sha256"
	"fmt"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"math/rand"
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
	CtrDir = "/var/lib/vessel/containers"
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
	cg.setPath(filepath.Join("vessel", c.digest)).
		setMemorySwapLimit(c.mem, c.swap).
		setCPULimit(c.cpus).
		setProcessLimit(c.pids)

	return cg.Load()
}

func (c *container) removeCGroups() error {
	cg := newCGroup()
	cg.setPath(filepath.Join("vessel", c.digest))
	return cg.Remove()
}

func (c *container) mountFromImage(img v1.Image) (func() error, error) {
	target := filepath.Join(CtrDir, c.digest, "mnt")
	layers, err := getImageLayers(img)
	if err != nil {
		return nil, errorWithMessage(err, "unable to get image layers")
	}
	return overlayFsMount(target, layers, false)
}
