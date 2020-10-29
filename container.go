package main

import (
	"crypto/md5"
	"math/rand"
	"time"
)

type containerOptions struct {
	name string
	hostname string
	mem  int
	swap int
	pid  int
	cpus float64
	detach bool
}

func (opt *containerOptions) setRandomName() {
	rand.Seed(time.Now().Unix())
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	opt.name = string(md5.New().Sum(randBuffer))
}