package main

import (
	"crypto/sha256"
	"math/rand"
	"syscall"
	"time"
)

type container struct {
	image    string
	digest   string
	name     string
	hostname string
	mem      int
	swap     int
	pids     int
	cpus     float64
	detach   bool
}

func (c *container) setDigest() {
	rand.Seed(time.Now().Unix())
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	sha := sha256.New().Sum(randBuffer)
	c.digest = string(sha)
}

func (c *container) setHostname() {
	if c.hostname == "" {
		c.hostname = c.image[:12]
	}
	Must(syscall.Sethostname([]byte(c.hostname)))
}

func (c *container) loadCGroups() (remover, error) {
	return newCGroup(c).Load()
}