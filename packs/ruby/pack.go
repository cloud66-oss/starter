package ruby

import "github.com/cloud66/starter/packs"

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

func (p *Pack) Name() string {
	return "ruby"
}

func (p *Pack) Detector() packs.Detector {
	return &Detector{PackElement: packs.PackElement{Pack: p}}
}

func (p *Pack) Analyze(rootDir string, environment string, shouldNotPrompt bool) error {
	var err error
	a := Analyzer{
		AnalyzerBase: packs.AnalyzerBase{
			PackElement:     packs.PackElement{Pack: p},
			RootDir:         rootDir,
			ShouldNotPrompt: shouldNotPrompt,
			Environment:     environment}}
	p.Analysis, err = a.Analyze()
	return err
}

func (p *Pack) WriteDockerfile(templateDir string, outputDir string, shouldNotPrompt bool) error {
	w := DockerfileWriter{
		packs.DockerfileWriterBase{
			PackElement: packs.PackElement{p},
			TemplateWriterBase: packs.TemplateWriterBase{
				TemplateDir:     templateDir,
				OutputDir:       outputDir,
				ShouldNotPrompt: shouldNotPrompt}}}
	return w.Write(p.Analysis.DockerfileContext)
}

func (p *Pack) WriteServiceYAML(templateDir string, outputDir string, shouldNotPrompt bool) error {
	w := ServiceYAMLWriter{
		packs.ServiceYAMLWriterBase{
			PackElement: packs.PackElement{p},
			TemplateWriterBase: packs.TemplateWriterBase{
				TemplateDir:     templateDir,
				OutputDir:       outputDir,
				ShouldNotPrompt: shouldNotPrompt}}}
	return w.Write(p.Analysis.ServiceYAMLContext)
}

func (p *Pack) GetMessages() []string {
	return p.Analysis.Messages.Items
}
