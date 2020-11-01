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

func NewTarGzFile(filepath string) (Extractor, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	return &TarGz{reader: reader}, nil
}

func NewTarGzExtractor(r io.Reader) Extractor {
	return &TarGz{reader: r}
}

func (t *TarGz) Extract(dst string) error {
	reader, err := gzip.NewReader(t.reader)
	if err != nil {
		return err
	}

	return NewTarExtractor(reader).Extract(dst)
}
