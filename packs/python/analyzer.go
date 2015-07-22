package python

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Analyzer struct {
	packs.AnalyzerBase
	RequirementsTxt string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.RequirementsTxt = filepath.Join(a.RootDir, "requirements.txt")
	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}

	packages := a.GuessPackages()
	version := a.FindVersion()
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()

	services, err := a.AnalyzeServices(a, envVars, gitBranch, gitURL, buildRoot)
	if err != nil {
		return nil, err
	}

	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:  a.GetPack().Name(),
			GitBranch: gitBranch,
			GitURL:    gitURL,
			Messages:  a.Messages},
		ServiceYAMLContext: &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs.Items}},
		DockerfileContext:  &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	return nil
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	return packages
}

func (a *Analyzer) FindVersion() string {
	hasFound, version := common.GetPythonVersion()
	return a.ConfirmVersion(hasFound, version, "latest")
}

func (a *Analyzer) FindDatabases() *common.Lister {
	dbs := common.NewLister()
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}
