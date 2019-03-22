package ruby

import (
	"fmt"
	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/packs"
	"os"
)

type Pack struct {
	packs.PackBase
	Analysis *Analysis
}

const (
	rubyRailsStencilTemplatePath = "https://raw.githubusercontent.com/cloud66/stencils-ruby-rails/{{.branch}}/" // this way we only have to add the filename. We should start by download the templates.json, do a couples of checks and after that download the stuff
)

func (p *Pack) Name() string {
	return "ruby"
}

func (p *Pack) LanguageVersion() string {
	return ""
}

func (p *Pack) FilesToBeAnalysed() []string {
	return []string{"Gemfile", "Procfile", "config/database.yml"}
}

func (p *Pack) Framework() string {
	return p.Analysis.Framework
}

func (p *Pack) FrameworkVersion() string {
	return p.Analysis.FrameworkVersion
}

func (p *Pack) GetSupportedLanguageVersions() []string {
	return nil
}

func (p *Pack) SetSupportedLanguageVersions(version []string) {

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

func (p *Pack) CreateSkycapFiles(outputDir string) error{
	var services = p.Analysis.ServiceYAMLContext.Services
	var envVars = make(map[string]string)
	for _, envVarArray := range services {
		for _, envs := range envVarArray.EnvVars {
			envVars[envs.Key] = envs.Value
		}
	}
	fmt.Printf("Env vars array: %v \n", envVars)

	//analyze the app, select the right template-repo, download the right stencils and helm releases, populate them with the rights values and create bundle file
	var templateRepository string = p.StencilRepositoryPath()
	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
	}else{
		//start download the template.json file
		fmt.Printf("template repo path: %s \n", templateRepository)
		getStencilTemplates(templateRepository)
		common.PrintlnL0("Now you can use the bundle file to create your formation with the following cx command:")
		common.PrintlnL0("magical cx command that will do everything")
	}






	return nil
}


func (p *Pack) GetMessages() []string {
	return p.Analysis.Messages.Items
}

func (p *Pack) GetDatabases() []string {
	return []string{}
}

func (p *Pack) GetStartCommands() []string {
	return p.Analysis.ListOfStartCommands
}

func (p *Pack) StencilRepositoryPath() string {
	return rubyRailsStencilTemplatePath
}

// downloading templates from github and putting them into homedir
func getStencilTemplates(repoPath string) error {
	tempDir := "./.skycap"
	common.PrintlnL0("Checking templates in %s", tempDir)

	//Create .bundle directory if it doesn't exist
	err := os.MkdirAll(tempDir, 0777)
	if err != nil {
		return err
	}

	//Download templates.json file
	manifest_path := repoPath+"templates.json"
	down_err := common.DownloadSingleFile(tempDir, common.DownloadFile{URL: manifest_path, Name: "templates.json"}, "master")
	if down_err != nil {
		return down_err
	}



	return nil
}