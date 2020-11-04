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
	containerPath       = "/var/run/vessel/containers"
	containerNetNsPath  = "/var/run/vessel/netns"
	containerConfigFile = "config.json"
	DigestStdLen        = 64
	MB                  = 1 << 20
)

type Container struct {
	Config *v1.Config
	Digest string
	RootFS string
	Pids   []int
	mem    int
	swap   int
	pids   int
	cpus   float64
}

/// NewContainer returns a new Container with a random digest.
func NewContainer() *Container {
	ctr := &Container{
		Config: new(v1.Config),
		Digest: randomHash(),
	}
	return ctr
}

// SetHostname sets Hostname for container.
//
// If Hostname was empty it uses the digest[:12]
func (c *Container) SetHostname() {
	if c.Config.Hostname == "" {
		c.Config.Hostname = c.Digest[:12]
	}
	syscall.Sethostname([]byte(c.Config.Hostname))
}

// Remove removes Container directory. It only works if all
// mount points have unmounted.
func (c *Container) Remove() error {
	if err := os.RemoveAll(filepath.Join(containerPath, c.Digest)); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(containerNetNsPath, c.Digest)); err != nil {
		return err
	}
	return c.removeCGroups()
}

// MountFromImage mounts filesystem for Container from an Image.
// It uses overlayFS for union mount of multiple layers.
func (c *Container) MountFromImage(img *image.Image) (filesystem.Unmounter, error) {
	target := filepath.Join(containerPath, c.Digest, "mnt")
	if err := os.MkdirAll(target, 0700); err != nil {
		return nil, errors.Wrapf(err, "can't create %s directory", target)
	}

	c.RootFS = target

	imgLayers, err := img.Layers()
	layers := make([]string, 0)
	for i := range imgLayers {
		digest, err := imgLayers[i].Digest()
		if err != nil {
			return nil, err
		}
		layers = append(layers, filepath.Join(image.LyrDir, digest.Hex))
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get image layers")
	}

	unmounter, err := filesystem.OverlayMount(target, layers, false)
	if err != nil {
		return unmounter, err
	}

	return unmounter, c.copyImageConfig(img)
}

// copyImageConfig copies image config into Container Directory.
func (c *Container) copyImageConfig(img v1.Image) error {
	file := filepath.Join(containerPath, c.Digest, containerConfigFile)
	data, err := img.RawConfigFile()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0655)
}

// LoadConfig loads and sets Container Config from its image config file.
func (c *Container) LoadConfig() error {
	filename := filepath.Join(containerPath, c.Digest, containerConfigFile)
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

// GetAllContainers returns slice of running Containers.
func GetAllContainers() ([]*Container, error) {
	all := make([]*Container, 0)

	list, err := ioutil.ReadDir(containerPath)
	if err != nil {
		return nil, err
	}
	for _, file := range list {
		if !file.IsDir() {
			continue
		}
		ctr, err := GetContainerByDigest(file.Name())
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		all = append(all, ctr)
	}

	return all, nil
}

// GetContainerByDigest returns a Container associate with a digest.
func GetContainerByDigest(digest string) (*Container, error) {
	ctrDigest := completeDigest(digest)
	if len(ctrDigest) != DigestStdLen {
		return nil, errors.Errorf("No such container: %s", digest)
	}

	config, err := getConfigByDigest(ctrDigest)
	if err != nil {
		return nil, err
	}
	pids, err := getPidsByDigest(ctrDigest)
	if err != nil {
		return nil, err
	}
	ctr := &Container{
		Config: config,
		RootFS: filepath.Join(containerPath, ctrDigest, "mnt"),
		Digest: ctrDigest,
		Pids:   pids,
	}
	return ctr, nil
}

func getConfigByDigest(digest string) (*v1.Config, error) {
	cfgPath := filepath.Join(containerPath, digest, containerConfigFile)
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
