package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func unTar(r io.Reader, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	hardLinks := make(map[string]string)
	tarReader := tar.NewReader(r)

loop:
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			break loop
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
				continue loop
			case err != nil:
				return err
			}

			if _, err = io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
		case tar.TypeLink:
			/* Store hardlinks for further finally*/
			link := filepath.Join(dst, header.Name)
			linkTarget := filepath.Join(dst, header.Linkname)
			hardLinks[link] = linkTarget
		case tar.TypeSymlink:
			linkPath := filepath.Join(dst, header.Name)
			if err := os.Symlink(header.Linkname, linkPath); err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
		}
	}

	for link, linkTarget := range hardLinks {
		if err := os.Link(link, linkTarget); err != nil {
			return err
		}
	}
	return nil
}
