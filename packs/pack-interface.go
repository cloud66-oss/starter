package packs

type Pack interface {
	Name() string
	Detector() Detector
	Analyze(rootDir string, environment string, shouldNotPrompt bool) error
	WriteDockerfile(templateDir string, outputDir string, shouldNotPrompt bool) error
	WriteServiceYAML(templateDir string, outputDir string, shouldNotPrompt bool) error
	WriteDockerComposeYAML(templateDir string, outputDir string, shouldNotPrompt bool) error
	GetMessages() []string
}
