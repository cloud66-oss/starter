package packs

type Detector interface {
	PackName() string
	Detect(rootDir string) bool
}
