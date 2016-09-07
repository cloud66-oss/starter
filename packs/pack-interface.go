package packs

type Pack interface {
	Name() string
	Framework() string
	FrameworkVersion() string
	FilesToBeAnalysed() [] string
	Detector() Detector
	Analyze(rootDir string, environment string, shouldNotPrompt bool, git_repo string, git_branch string) error
	WriteDockerfile(templateDir string, outputDir string, shouldNotPrompt bool) error
	WriteServiceYAML(templateDir string, outputDir string, shouldNotPrompt bool) error
	WriteDockerComposeYAML(templateDir string, outputDir string, shouldNotPrompt bool) error
	GetMessages() []string
}
