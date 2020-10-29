package main

import (
	"crypto/md5"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"math/rand"
	"time"
)

type containerOptions struct {
	image    v1.Image
	name     string
	hostname string
	mem      int
	swap     int
	pid      int
	cpus     float64
	detach   bool
}

func (opt *containerOptions) setRandomName() {
	rand.Seed(time.Now().Unix())
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	opt.name = string(md5.New().Sum(randBuffer))
}
