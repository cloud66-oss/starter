package packs

type Detector interface {
	Name() string
	Detect(root string) bool
	Analyzer(root string, environment string) Pack
}
