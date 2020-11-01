package archive

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Tarball struct {
	reader io.Reader
}

func NewTarFile(filepath string) (Extractor, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	return &Tarball{reader: reader}, nil
}

func NewTarExtractor(r io.Reader) Extractor {
	return &Tarball{reader: r}
}

func (t *Tarball) Extract(dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	tarReader := tar.NewReader(t.reader)

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		path := filepath.Join(dst, header.Name)
		info := header.FileInfo()

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			switch {
			case os.IsExist(err):
				continue
			case err != nil:
				return err
			}

			if _, err = io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
		case tar.TypeLink:
			link := filepath.Join(dst, header.Name)
			linkTarget := filepath.Join(dst, header.Linkname)
			// lazy link creation. just to make sure all files are available
			defer os.Link(link, linkTarget)
		case tar.TypeSymlink:
			linkPath := filepath.Join(dst, header.Name)
			if err := os.Symlink(header.Linkname, linkPath); err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
		}
	}
}
