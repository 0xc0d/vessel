package archive

type Extractor interface {
	Extract(dst string) error
}
