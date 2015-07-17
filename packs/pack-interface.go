package packs

type Pack interface {
	Name() string
	Detector() Detector
	Analyze(rootDir string, environment string) error
	WriteDockerfile(string, string, bool) error
	WriteServiceYAML(string, string, bool) error
	GetMessages() []string
}
