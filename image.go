package main

import (
	"compress/gzip"
	"encoding/json"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TmpDir   = "/home/aggy/var/lib/vessel/temp"
	ImgDir   = "/home/aggy/var/lib/vessel/images"
	tarExt   = ".tar"
	tarGzExt = ".tar.gz"
	manifest = "manifest.json"
)

func init() {
	Must(os.MkdirAll(TmpDir, 0755))
	Must(os.MkdirAll(ImgDir, 0755))
}

type Manifest struct {
	Config string
	Layers []string
}

func getImage(src string) (v1.Image, error) {
	img, err := crane.Pull(src)
	if err != nil {
		return nil, err
	}
	if !imageExists(img) {
		if err := imageDownload(img); err != nil {
			return nil, err
		}
	}

	return img, nil
}

// imageExists checks for image existence in local storage
func imageExists(img v1.Image) bool {
	imageDigest, err := img.Digest()
	CheckErr(err)
	files, err := ioutil.ReadDir(ImgDir)
	CheckErr(err)
	for _, file := range files {
		if file.IsDir() && imageDigest.Hex == file.Name() {
			return true
		}
	}
	return false
}

// imageDownload downloads tar format of image in TmpDir.
// It then extract the image and copy it into ImgDir
func imageDownload(img v1.Image) error {
	imageDigest, err := img.Digest()
	CheckErr(err)

	imageTmpPath := filepath.Join(TmpDir, imageDigest.Hex)+tarExt
	if err := crane.Save(img, imageDigest.Hex, imageTmpPath); err != nil {
		return err
	}
	if err := imageExtract(img); err != nil {
		return err
	}

	return imageRemove(img)
}

// imageExtract extracts the downloaded image
func imageExtract(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageTempDir := filepath.Join(TmpDir, imageDigest.Hex)
	imageFile := imageTempDir + tarExt

	file, err := os.Open(imageFile)
	if err != nil {
		return err
	}

	if err := unTar(file, imageTempDir); err != nil {
		return err
	}

	// Extract image layers
	if err := imageLayerExtract(img); err != nil {
		return err
	}

	// Move manifest.json to ImgDir
	if err := os.Rename(filepath.Join(imageTempDir, manifest),
		filepath.Join(ImgDir, imageDigest.Hex, manifest)); err != nil {
		return err
	}

	// Move config file to ImgDir
	newManifest, err := parseImageManifest(imageDigest.Hex)
	if err != nil {
		return err
	}
	return os.Rename(filepath.Join(imageTempDir, newManifest.Config),
		filepath.Join(ImgDir, imageDigest.Hex, newManifest.Config))
}

// imageLayerExtract tar Gzip layer files in image
func imageLayerExtract(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	manifest, err := img.Manifest()
	if err != nil {
		return err
	}

	ch := make(chan error, len(manifest.Layers))

	for _, layer := range manifest.Layers {
		layer := layer
		layerFilepath := filepath.Join(TmpDir, imageDigest.Hex, layer.Digest.Hex) + tarGzExt
		file, err := os.Open(layerFilepath)
		if err != nil {
			return err
		}
		reader, err := gzip.NewReader(file)
		go func() {
			defer file.Close()
			ch <- unTar(reader, filepath.Join(ImgDir, imageDigest.Hex, layer.Digest.Hex))
		}()
	}

	for i := 0; i < len(manifest.Layers); i++ {
		err := <-ch
		if err != nil {
			return err
		}
	}

	return nil
}

func imageRemove(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageTempDir := filepath.Join(TmpDir, imageDigest.Hex)
	imageFile := imageTempDir + tarExt
	if err := os.RemoveAll(imageTempDir); err != nil {
		return err
	}
	return  os.RemoveAll(imageFile)
}


func getImageLayers(img v1.Image) ([]string, error) {
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

// parseImageManifest parses manifest.json file of an image
func parseImageManifest(src string) (Manifest, error) {
	manifestFile, err := ioutil.ReadFile(filepath.Join(ImgDir, src, manifest))
	if err != nil {
		return Manifest{}, err
	}

	var ml []Manifest
	err = json.Unmarshal(manifestFile, &ml)
	return ml[0], err
}

// parseImageConfig parses config file of an image
func parseImageConfig(src string) (v1.Config, error) {
	var config v1.Config
	manifest, err := parseImageManifest(src)
	if err != nil {
		return config, err
	}
	file, err := os.Open(filepath.Join(ImgDir, src, manifest.Config))
	if err != nil {
		return config, err
	}

	configFile, err := v1.ParseConfigFile(file)
	if err != nil {
		return config, err
	}
	return configFile.Config, err
}
