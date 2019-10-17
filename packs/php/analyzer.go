package php

import (
	"path/filepath"

	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/packs"
)

type Analyzer struct {
	packs.AnalyzerBase
	ComposerJSON string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.ComposerJSON = filepath.Join(a.RootDir, "composer.json")
	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}
	version := a.FindVersion()
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()
	packages := a.GuessPackages()
	framework := a.GuessFramework()
	a.CheckNotSupportedPackages(packages)

	services, err := a.AnalyzeServices(a, envVars, gitBranch, gitURL, buildRoot)

	// inject all the services with the databases used in the infrastructure
	for _, service := range services {
		service.Databases = dbs
	}

	if err != nil {
		return nil, err
	}

	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:  a.GetPack().Name(),
			GitBranch: gitBranch,
			GitURL:    gitURL,
			Framework: framework,
			Messages:  a.Messages},
		DockerComposeYAMLContext: &DockerComposeYAMLContext{packs.DockerComposeYAMLContextBase{Services: services, Dbs: dbs}},
		ServiceYAMLContext:       &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs}},
		DockerfileContext:        &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	if len(*services) == 0 {
		var service *common.Service
		service = &common.Service{Name: "web"}
		service.Ports = []*common.PortMapping{common.NewPortMapping()}
		service.Command = "apache2-foreground"
		service.Ports[0].Container = "80"
		*services = append(*services, service)
	}
	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	hasFound, _ := common.GetDependencyVersion(a.ComposerJSON, pack)
	return hasFound
}

func (a *Analyzer) GuessFramework() string {
	if runsLaravel, _ := common.GetFramework(a.ComposerJSON, "laravel/framework"); runsLaravel {
		return "laravel"
	}
	return ""
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	return packages
}

func (a *Analyzer) FindVersion() string {
	foundNode, phpVersion := common.GetPHPVersion(a.ComposerJSON)
	return a.ConfirmVersion(foundNode, phpVersion, "latest")
}

func (a *Analyzer) FindDatabases() []common.Database {
	dbs := []common.Database{}
	if hasMysql, _ := common.GetPHPDatabase(a.ComposerJSON, "mysql"); hasMysql {
		dbs = append(dbs, common.Database{Name: "mysql", DockerImage: "mysql"})
	}
	//if hasPostgres, _ := common.GetPHPDatabase(a.ComposerJSON, "pgsql"); hasPostgres {
	//	dbs = append(dbs, common.Database{Name: "postgres", DockerImage: "postgres"})
	//}
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}
