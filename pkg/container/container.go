package container

import (
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
	CtrDir       = "/var/lib/vessel/containers"
	CtrCfg       = "config"
	DigestStdLen = 64
	MB           = 1 << 20
)

type Container struct {
	Config *v1.Config
	Digest string
	RootFS string
	Name   string
	Pid    int
	mem    int
	swap   int
	pids   int
	cpus   float64
}

func NewContainer() *Container {
	ctr := &Container{
		Config: new(v1.Config),
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
	file := filepath.Join(CtrDir, c.Digest, CtrCfg)
	data, err := img.RawConfigFile()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0655)
}

func (c *Container) LoadConfig() error {
	filename := filepath.Join(CtrDir, c.Digest, CtrCfg)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	configFile, err := v1.ParseConfigFile(file)
	if err != nil {
		return err
	}
	c.Config = configFile.Config.DeepCopy()
	return nil
}

func GetAllContainer() ([]*Container, error) {
	all := make([]*Container, 0)

	list, err := ioutil.ReadDir(CtrDir)
	if err != nil {
		return nil, err
	}
	for _, file := range list {
		if !file.IsDir() {
			continue
		}
		ctr, err := GetContainerByDigest(file.Name())
		if err != nil {
			return nil, err
		}
		all = append(all, ctr)
	}

	return all, nil
}

func GetContainerByDigest(digest string) (*Container, error) {
	ctrDigest := completeDigest(digest)
	if len(ctrDigest) != DigestStdLen {
		return nil, errors.Errorf("No such container: %s", digest)
	}

	config, err := GetConfigByDigest(ctrDigest)
	if err != nil {
		return nil, err
	}
	pid, err := GetPidByDigest(ctrDigest)
	if err != nil {
		return nil, err
	}
	ctr := &Container{
		Config: config,
		RootFS: filepath.Join(CtrDir, ctrDigest, "mnt"),
		Digest: ctrDigest,
		Pid:    pid,
	}
	return ctr, nil
}

func GetConfigByDigest(digest string) (*v1.Config, error) {
	cfgPath := filepath.Join(CtrDir, digest, CtrCfg)
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()
	cfg, err := v1.ParseConfigFile(cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg.Config.DeepCopy(), nil
}
