package main

import (
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

const (
	BasePath = "/var/lib/vessel"
	TmpPath  = "/var/lib/vessel/tmp"
	ImgPath  = "/var/lib/vessel/images"
	CtrPath  = "/var/lib/vessel/containers"
)

// newImage returns v1.Image for the given image source
func newImage(src string) (v1.Image, error) {
	return crane.Pull(src)
}