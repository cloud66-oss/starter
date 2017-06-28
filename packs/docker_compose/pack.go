package docker_compose

import ()
import (
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/common"
)

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

func (p *Pack) Name() string {
	return "docker-compose"
}

func (p *Pack) LanguageVersion() string {
	return p.Analysis.LanguageVersion
}

func (p *Pack) FilesToBeAnalysed() [] string {
	return []string{"docker-compose.yml"}
}

func (p *Pack) Framework() string {
	return ""
}

func (p *Pack) FrameworkVersion() string {
	return ""
}

func (p *Pack) GetSupportedLanguageVersions() []string {
	return nil
}

func (p *Pack) SetSupportedLanguageVersions(versions []string) {

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

	common.PrintlnWarning("You can not generate a Dockerfile using this pack. Nothing to do.")
	return nil
}

func (p *Pack) WriteServiceYAML(templateDir string, outputDir string, shouldPrompt bool) error {

	err := Transformer(outputDir+"/docker-compose.yml", outputDir+"/service.yml", p.Analysis.GitURL, p.Analysis.GitBranch, shouldPrompt)

	CheckError(err)

	return nil
}

func (p *Pack) WriteDockerComposeYAML(templateDir string, outputDir string, shouldPrompt bool) error {
	common.PrintlnWarning("There is already an existing docker-compose.yml. Nothing to do.")
	return nil
}

func (p *Pack) GetMessages() []string {
	return []string{}
}

func (p *Pack) GetDatabases() []string {
	return []string{}
}

func (p *Pack) GetStartCommands() []string {
	return []string{}
}
