package main

import (
	"compress/gzip"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TmpPath  = "/home/aggy/var/lib/vessel/tmp"
	ImgPath  = "/home/aggy/var/lib/vessel/images"
	tarExt   = ".tar"
	tarGzExt = ".tar.gz"
)

func init() {
	Must(os.MkdirAll(TmpPath, 0755))
	Must(os.MkdirAll(ImgPath, 0755))
}

// newImage returns v1.Image for the given image source
func newImage(src string) (v1.Image, error) {
	return crane.Pull(src)
}

// imageExists checks for image existence in local storage
func imageExists(img v1.Image) bool {
	imageDigest, err := img.Digest()
	CheckErr(err)
	files, err := ioutil.ReadDir(ImgPath)
	CheckErr(err)
	for _, file := range files {
		if file.IsDir() && imageDigest.Hex == file.Name() {
			return true
		}
	}
	return false
}

// imageDownload downloads tar format of image in TmpPath.
// It then extract the image and copy it into ImgPath
func imageDownload(img v1.Image) error {
	imageDigest, err := img.Digest()
	CheckErr(err)

	err = crane.Save(img, imageDigest.Hex, filepath.Join(TmpPath, imageDigest.Hex)+tarExt)
	if err != nil {
		return err
	}

	return imageExtract(img)
}

// imageExtract extracts the downloaded image
func imageExtract(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageDir := filepath.Join(TmpPath, imageDigest.Hex)
	defer os.RemoveAll(imageDir)
	imageFilepath := imageDir + tarExt
	defer os.RemoveAll(imageFilepath)

	file, err := os.Open(imageFilepath)
	if err != nil {
		return err
	}

	if err := unTar(file, imageDir); err != nil {
		return err
	}

	if err := layerExtract(img); err != nil {
		return err
	}

	return nil
}

// layerExtract tar Gzip layer files in image
func layerExtract(img v1.Image) error {
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
		layerFilepath := filepath.Join(TmpPath, imageDigest.Hex, layer.Digest.Hex) + tarGzExt
		file, err := os.Open(layerFilepath)
		if err != nil {
			return err
		}
		reader, err := gzip.NewReader(file)
		go func() {
			defer file.Close()
			ch <- unTar(reader, filepath.Join(ImgPath, imageDigest.Hex, layer.Digest.Hex))
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
