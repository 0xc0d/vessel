package archive

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

type TarGz struct {
	reader io.Reader
}

// NewTarGzFile creates a Gziped tarball for the given filename.
func NewTarGzFile(filename string) (Extractor, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	return &TarGz{reader: reader}, nil
}

// NewTarGz creates a Gziped tarball for the give Reader.
func NewTarGz(r io.Reader) Extractor {
	return &TarGz{reader: r}
}


// Extract extracts a Gziped tarball into dst.
func (t *TarGz) Extract(dst string) error {
	reader, err := gzip.NewReader(t.reader)
	if err != nil {
		return err
	}

	return NewTar(reader).Extract(dst)
}
