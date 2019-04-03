package service_yml_to_kubes

import (
	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/definitions/service-yml"
	"github.com/cloud66-oss/starter/packs"
	"github.com/cloud66-oss/starter/transform"
)

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

const (
	StencilTemplatePath = "" //still not implemented
	templateRepositoryBranch = ""
)

func (p *Pack) Name() string {
	return "service.yml"
}

func (p *Pack) LanguageVersion() string {
	return p.Analysis.LanguageVersion
}

func (p *Pack) FilesToBeAnalysed() []string {
	return []string{"service.yml"}
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

	//common.PrintlnWarning("You can not generate a Dockerfile using this pack. Nothing to do.")
	return nil
}

func (p *Pack) WriteKubesConfig(outputDir string, shouldPrompt bool) error {

	serviceYmlBase := service_yml.ServiceYml{}
	serviceYmlBase.UnmarshalFromFile(outputDir + "/service.yml")

	s := transform.ServiceYmlTransformer{Base: serviceYmlBase}

	kubernetes := s.ToKubernetes()
	kubernetes.MarshalToFile(outputDir + "/kubernetes.yml")

	return nil
}

func (p *Pack) WriteDockerComposeYAML(templateDir string, outputDir string, shouldPrompt bool) error {
	common.PrintlnWarning("You can not generate a docker-compose.yml using this pack. Nothing to do.")
	return nil
}
func (p *Pack) WriteServiceYAML(templateDir string, outputDir string, shouldPrompt bool) error {
	common.PrintlnWarning("There is already an existing service.yml. Nothing to do.")
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

func (p *Pack) StencilRepositoryPath() (string, string) {
	return StencilTemplatePath, templateRepositoryBranch
}

func (p *Pack) CreateSkycapFiles(outputDir string, templateDir string) error{
	common.PrintlnWarning("You can not generate the Skycap configuration files using this pack. Nothing to do.")
	return nil
}