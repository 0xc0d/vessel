package image

import (
	"encoding/json"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io/ioutil"
	"os"
)

const (
	ImgDir   = "/var/lib/vessel/images"
	RepoFile = "/var/lib/vessel/images/repositories.json"
	LyrDir   = "/var/lib/vessel/images/layers"
)

type Image struct {
	v1.Image
	Registry   string
	Repository string
	Name       string
	Tag        string
}

func NewImage(src string) (*Image, error) {
	tag, err := name.NewTag(src)
	img, err := crane.Pull(tag.Name())
	if err != nil {
		return nil, err
	}

	newImage := &Image{
		Image:      img,
		Registry:   tag.RegistryStr(),
		Repository: tag.RepositoryStr(),
		Name:       tag.Name(),
		Tag:        tag.TagStr(),
	}
	return newImage, nil
}

// Exists checks for image existence in local storage.
func (i *Image) Exists() (bool, error) {
	repos, err := GetAll()
	if err != nil {
		return false, err
	}
	for _, repo := range repos {
		for repoName, _ := range repo {
			if repoName == i.Name {
				return true, nil
			}
		}
	}
	return false, nil
}

func GetAll() (Repositories, error) {
	repos := make(Repositories)

	data, err := ioutil.ReadFile(RepoFile)
	if os.IsNotExist(err) {
		return repos, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}
