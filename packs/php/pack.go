package php

import "github.com/cloud66/starter/packs"

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

func (p *Pack) Name() string {
	return "php"
}

func (p *Pack) FilesToBeAnalysed() []string {
	return []string{ "composer.json" }
}

func (p *Pack) LanguageVersion() string {
	return "x"
}


func (p *Pack) Framework() string {
	return p.Analysis.Framework
}

func (p *Pack) FrameworkVersion() string {
	return p.Analysis.FrameworkVersion
}

func (p *Pack) Detector() packs.Detector {
	return &Detector{PackElement: packs.PackElement{Pack: p}}
}

func (p *Pack) Analyze(rootDir string, environment string, shouldPrompt bool, git_repo string, git_branch string) error {
	var err error
	a := Analyzer{
		AnalyzerBase: packs.AnalyzerBase{
			PackElement:  packs.PackElement{Pack: p},
			RootDir:      rootDir,
			ShouldPrompt: shouldPrompt,
			Environment:  environment}}
	p.Analysis, err = a.Analyze()
	return err
}

func (p *Pack) WriteDockerfile(templateDir string, outputDir string, shouldPrompt bool) error {
	w := DockerfileWriter{
		packs.DockerfileWriterBase{
			PackElement: packs.PackElement{Pack: p},
			TemplateWriterBase: packs.TemplateWriterBase{
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				ShouldPrompt: shouldPrompt}}}
	return w.Write(p.Analysis.DockerfileContext)
}

func (p *Pack) WriteServiceYAML(templateDir string, outputDir string, shouldPrompt bool) error {
	w := ServiceYAMLWriter{
		packs.ServiceYAMLWriterBase{
			PackElement: packs.PackElement{Pack: p},
			TemplateWriterBase: packs.TemplateWriterBase{
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				ShouldPrompt: shouldPrompt}}}
	return w.Write(p.Analysis.ServiceYAMLContext)
}

func (p *Pack) WriteDockerComposeYAML(templateDir string, outputDir string, shouldPrompt bool) error {
	w := DockerComposeYAMLWriter{
		packs.DockerComposeYAMLWriterBase{
			PackElement: packs.PackElement{Pack: p},
			TemplateWriterBase: packs.TemplateWriterBase{
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				ShouldPrompt: shouldPrompt}}}
	return w.Write(p.Analysis.DockerComposeYAMLContext)
}

func (p *Pack) GetMessages() []string {
	return p.Analysis.Messages.Items
}

func (p *Pack) GetDatabases() []string {
	return []string {}
}

func (p *Pack) GetStartCommands() []string {
	return  p.Analysis.ListOfStartCommands
}
