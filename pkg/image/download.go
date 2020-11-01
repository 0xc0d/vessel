package image

import (
	"github.com/0xc0d/vessel/pkg/archive"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1"
	"os"
	"path/filepath"
)

// download downloads tar format of image in TmpDir.
// It then extract the image and copy it into ImgDir
func Download(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageTmpPath := filepath.Join(TmpDir, imageDigest.Hex) + tarExt
	if err := crane.Save(img, imageDigest.Hex, imageTmpPath); err != nil {
		return err
	}
	if err := extract(img); err != nil {
		return err
	}

	return remove(img)
}

// extract extracts the downloaded image
func extract(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageTempDir := filepath.Join(TmpDir, imageDigest.Hex)
	imageFile := imageTempDir + tarExt

	tar, err := archive.NewTarFile(imageFile)
	if err != nil {
		return err
	}
	if err := tar.Extract(imageTempDir); err != nil {
		return err
	}

	// extract image layers
	if err := layerExtract(img); err != nil {
		return err
	}

	return nil
}

// layerExtract extracts tar.gz layerFiles inside image
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
		layerFilepath := filepath.Join(TmpDir, imageDigest.Hex, layer.Digest.Hex) + tarGzExt
		tarGz, err := archive.NewTarGzFile(layerFilepath)
		if err != nil {
			return err
		}
		dst := filepath.Join(ImgDir, imageDigest.Hex, layer.Digest.Hex)
		go func() {
			ch <- tarGz.Extract(dst)
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

func remove(img v1.Image) error {
	imageDigest, err := img.Digest()
	if err != nil {
		return err
	}

	imageTempDir := filepath.Join(TmpDir, imageDigest.Hex)
	imageFile := imageTempDir + tarExt
	if err := os.RemoveAll(imageTempDir); err != nil {
		return err
	}
	return os.RemoveAll(imageFile)
}