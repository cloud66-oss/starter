package node

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Analyzer struct {
	packs.AnalyzerBase
	PackageJSON string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.PackageJSON = filepath.Join(a.RootDir, "package.json")
	gitURL, gitBranch, buildRoot, _ := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}
	version := a.FindVersion()
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()
	packages := a.GuessPackages()
	a.CheckNotSupportedPackages(packages)

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
		DockerComposeYAMLContext: &DockerComposeYAMLContext{packs.DockerComposeYAMLContextBase{Services: services, Dbs: dbs.Items}},
		ServiceYAMLContext: &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs.Items}},
		DockerfileContext:  &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	//has a procfile
	if len(*services) == 1 {
		for _, service := range *services {
			service.Ports = []*common.PortMapping{common.NewPortMapping()}
			service.Ports[0].Container = "3000"
		}
	}


	//no procfile
	if len(*services) == 0 {
		var service *common.Service
		service = &common.Service{Name: "web"}
		service.Ports = []*common.PortMapping{common.NewPortMapping()}
		service.Command = "node index.js"
		service.Ports[0].Container = "3000"
		*services = append(*services, service)
	}

	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	hasFound, _ := common.GetDependencyVersion(a.PackageJSON, pack)
	return hasFound
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()

	if runsExpress, _ := common.GetDependencyVersion(a.PackageJSON, "express"); runsExpress {
		common.PrintlnL2("Found Express")
	}
	if hasScript, script := common.GetScriptsStart(a.PackageJSON); hasScript {
		common.PrintlnL2("Found Script: %s", script)
	}

	return packages
}

func (a *Analyzer) FindVersion() string {
	foundNode, nodeVersion := common.GetNodeVersion(a.PackageJSON)
	return a.ConfirmVersion(foundNode, nodeVersion, "latest")
}

func (a *Analyzer) FindDatabases() *common.Lister {
	dbs := common.NewLister()
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}
