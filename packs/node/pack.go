package node

import "github.com/cloud66-oss/starter/packs"
import "github.com/cloud66-oss/starter/common"

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

const (
	nodeExpressStencilTemplatePath = "" // this way we only have to add the filename. We should start by download the templates.json, do a couples of checks and after that download the stuff
	templateRepositoryBranch = "master"
)

func (p *Pack) Name() string {
	return "node"
}

func (p *Pack) LanguageVersion() string {
	return p.Analysis.LanguageVersion
}

func (p *Pack) FilesToBeAnalysed() []string {
	return []string{"package.json", "Procfile", ".meteor/release"}
}

func (p *Pack) Framework() string {
	return p.Analysis.Framework
}

func (p *Pack) FrameworkVersion() string {
	return p.Analysis.FrameworkVersion
}

func (p *Pack) GetSupportedLanguageVersions() []string {
	if p.Analysis != nil {
		return p.Analysis.SupportedLanguageVersions
	} else {
		return common.GetAllowedNodeVersions()
	}

}

func (p *Pack) SetSupportedLanguageVersions(versions []string) {
	common.SetAllowedNodeVersions(versions)
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
			GitURL:       git_repo,
			GitBranch:    git_branch,
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
	p.Analysis.DockerfileContext.FrameworkVersion = p.Analysis.FrameworkVersion
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

func (p *Pack) WriteKubesConfig(outputDir string, shouldPrompt bool) error {
	common.PrintlnWarning("You can not generate a Kubernetes configuration file using this pack. Nothing to do.")
	return nil
}

func (p *Pack) GetMessages() []string {
	return p.Analysis.Messages.Items
}

func (p *Pack) GetDatabases() []string {
	return p.Analysis.Databases
}

func (p *Pack) GetStartCommands() []string {
	return p.Analysis.ListOfStartCommands
}

func (p *Pack) StencilRepositoryPath() (string, string) {
	return nodeExpressStencilTemplatePath, templateRepositoryBranch
}

func (p *Pack) CreateSkycapFiles(outputDir string, templateDir string) error{
	common.PrintlnWarning("You can not generate the Skycap configuration files using this pack. Nothing to do.")
	return nil
}