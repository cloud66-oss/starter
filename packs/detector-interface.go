package packs

type Detector interface {
	GetPack() Pack
	Detect(rootDir string) bool
}
