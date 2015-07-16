package node

import "github.com/cloud66/starter/packs"

type Pack struct {
	packs.PackBase
}

func (p *Pack) Name() string {
	return "node"
}

func (p *Pack) Detector() packs.Detector {
	return &Detector{PackElement: packs.PackElement{Pack: p}}
}

func (p *Pack) Analyzer(rootDir string, environment string) packs.Analyzer {
	return &Analyzer{
		AnalyzerBase: packs.AnalyzerBase{PackElement: packs.PackElement{Pack: p},
			RootDir:     rootDir,
			Environment: environment}}
}

func (p *Pack) DockerfileWriter(templateDir string, outputDir string, shouldOverwrite bool) packs.DockerfileWriterBase {
	d := packs.DockerfileWriterBase{TemplateDir: templateDir, OutputDir: outputDir, ShouldOverwrite: shouldOverwrite}
	d.Pack = &Pack{}
	return d
}
