package container

import (
	"github.com/0xc0d/vessel/pkg/cgroups"
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/0xc0d/vessel/pkg/image"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

const (
	CtrDir = "/var/lib/vessel/containers"
	MB     = 1 << 20
)

type Container struct {
	Config v1.Config
	Digest string
	RootFS string
	Name   string
	mem    int
	swap   int
	pids   int
	cpus   float64
	detach bool
}

func NewContainer() *Container {
	ctr := &Container{
		Digest: randomHash(),
	}
	return ctr
}

func (c *Container) SetHostname() {
	if c.Config.Hostname == "" {
		c.Config.Hostname = c.Digest[:12]
	}
	syscall.Sethostname([]byte(c.Config.Hostname))
}

func (c *Container) LoadCGroups() error {
	cg := cgroups.NewCGroup()
	cg.SetPath(filepath.Join("vessel", c.Digest)).
		SetMemorySwapLimit(c.mem*MB, c.swap*MB).
		SetCPULimit(c.cpus).
		SetProcessLimit(c.pids)

	return cg.Load()
}

func (c *Container) RemoveCGroups() error {
	cg := cgroups.NewCGroup()
	cg.SetPath(filepath.Join("vessel", c.Digest))
	return cg.Remove()
}

func (c *Container) Remove() error {
	return os.RemoveAll(filepath.Join(CtrDir, c.Digest))
}

func (c *Container) MountFromImage(img v1.Image) (filesystem.Unmounter, error) {
	target := filepath.Join(CtrDir, c.Digest, "mnt")
	c.RootFS = target

	layers, err := image.GetLayers(img)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get image layers")
	}

	unmounter, err := filesystem.OverlayMount(target, layers, false)
	if err != nil {
		return unmounter, err
	}

	return unmounter, c.copyImageConfig(img)
}

func (c *Container) copyImageConfig(img v1.Image) error {
	file := filepath.Join(CtrDir, c.Digest, "config")
	data, err := img.RawConfigFile()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0655)
}

func (c *Container) LoadConfig() error {
	filename := filepath.Join(CtrDir, c.Digest, "config")
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	configFile, err := v1.ParseConfigFile(file)
	if err != nil {
		return err
	}
	c.Config = configFile.Config
	return nil
}
