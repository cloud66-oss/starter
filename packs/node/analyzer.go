package node

import (
	"fmt"
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
	var service *common.Service
	for _, s := range *services {
		if s.Name == "web" || s.Name == "custom_web" {
			service = s
			break
		}
	}
	if service == nil {
		service = &common.Service{Name: "web"}
		*services = append(*services, service)
	}
	return nil
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()

	if runsExpress, _ := common.GetDependencyVersion(a.PackageJSON, "express"); runsExpress {
		fmt.Println(common.MsgL2, "----> Found Express", common.MsgReset)
	}
	if hasScript, script := common.GetScriptsStart(a.PackageJSON); hasScript {
		fmt.Println(common.MsgL2, "----> Found Script:", script, common.MsgReset)
	}

	return packages
}

func (a *Analyzer) FindVersion() string {
	foundNode, nodeVersion := common.GetNodeVersion(a.PackageJSON)

	if foundNode {
		return fmt.Sprintf("%s-onbuild", nodeVersion)
	} else {
		nodeVersion = common.AskUser("Can't find Node version from package.json:", "default")
		if nodeVersion == "default" {
			return a.defaultVersion()
		} else {
			return fmt.Sprintf("%s-onbuild", nodeVersion)
		}
	}
}

func (a *Analyzer) FindDatabases() *common.Lister {
	dbs := common.NewLister()
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}

func (a *Analyzer) defaultVersion() string {
	return "onbuild"
}
