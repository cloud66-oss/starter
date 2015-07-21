package packs

type Pack interface {
	Name() string
	Detector() Detector
	Analyze(rootDir string, environment string) error
	WriteDockerfile(string, string) error
	WriteServiceYAML(string, string) error
	GetMessages() []string
}
