package packs

type Detector interface {
	Name() string
	Detect(rootDir string) bool
	Analyzer(rootDir string, environment string) Analyzer
}
