package packs

type Pack interface {
	Name() string
	Detector() Detector
	Analyzer(rootDir string, environment string) Analyzer
}

type PackElement struct {
	Pack Pack
}

func (e *PackElement) GetPack() Pack {
	return e.Pack
}
