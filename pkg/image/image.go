package image

import (
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TmpDir   = "/var/lib/vessel/temp"
	ImgDir   = "/var/lib/vessel/images"
	tarExt   = ".tar"
	tarGzExt = ".tar.gz"
)

func init() {
	os.MkdirAll(TmpDir, 0755)
	os.MkdirAll(ImgDir, 0755)
}

func NewImage(name string) (v1.Image, error) {
	img, err := crane.Pull(name)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GetLayers(img v1.Image) ([]string, error) {
	var layers []string
	imageDigest, err := img.Digest()
	if err != nil {
		return layers, err
	}
	imgManifest, err := img.Manifest()
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join(ImgDir, imageDigest.Hex)
	for _, layer := range imgManifest.Layers {
		layers = append(layers, filepath.Join(imagePath, layer.Digest.Hex))
	}
	return layers, nil
}

// exists checks for image existence in local storage
func Exists(img v1.Image) bool {
	imageDigest, err := img.Digest()

	files, err := ioutil.ReadDir(ImgDir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.IsDir() && imageDigest.Hex == file.Name() {
			return true
		}
	}
	return false
}
