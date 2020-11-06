package image

import (
	"encoding/json"
	"github.com/0xc0d/vessel/pkg/archive"
	"os"
	"path/filepath"
)

type Repository map[string]string
type Repositories map[string]Repository

// Download downloads image's layers.
func (i *Image) Download() error {
	layers, err := i.Layers()
	if err != nil {
		return err
	}
	for _, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			return err
		}
		rc, err := layer.Uncompressed()
		if err != nil {
			return err
		}

		tarball := archive.NewTar(rc)
		err = tarball.Extract(filepath.Join(LyrDir, digest.Hex))
		rc.Close()
		if err != nil {
			return err
		}
	}
	return i.addToRepositories()
}

// addToRepositories adds a image to repository.
//
// after downloading image's layer, call this function to prevent
// further duplicate downloads.
func (i *Image) addToRepositories() error {
	file, err := os.OpenFile(RepoFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	repos := make(Repositories)
	decoder := json.NewDecoder(file)
	decoder.Decode(&repos)

	// truncate file to overwrite
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	digest, err := i.Digest()
	if err != nil {
		return err
	}
	repos[i.Repository] = Repository{i.Name: digest.String()}
	encoder := json.NewEncoder(file)
	return  encoder.Encode(repos)
}