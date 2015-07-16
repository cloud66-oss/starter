package packs

type Pack interface {
	Name() string
	Detector() Detector
	Analyzer(rootDir string, environment string) Analyzer
	DockerfileWriter(string, string, bool) DockerfileWriterBase
	ServiceYAMLWriter(string, string, bool) ServiceYAMLWriterBase
}

type PackBase struct {
}

func (p *PackBase) ServiceYAMLWriter(templateDir string, outputDir string, shouldOverwrite bool) ServiceYAMLWriterBase {
	return ServiceYAMLWriterBase{TemplateDir: templateDir, OutputDir: outputDir, ShouldOverwrite: shouldOverwrite}
}

type PackElement struct {
	Pack Pack
}

func (e *PackElement) GetPack() Pack {
	return e.Pack
}
