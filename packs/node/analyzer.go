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

type Analysis struct {
	PackName string

	GitURL    string
	GitBranch string

	ServiceYAMLContext *ServiceYAMLContext
	DockerfileContext  *DockerfileContext

	Messages common.Lister
}

func (a *Analyzer) Init() error {
	a.PackageJSON = filepath.Join(a.GetRootDir(), "package.json")
	return nil
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	err := a.Init()
	if err != nil {
		fmt.Printf("%s Failed to initialize analyzer due to %s\n", common.MsgError, err.Error())
		return nil, err
	}

	gitURL := common.LocalGitBranch(a.GetRootDir())
	gitBranch := common.RemoteGitUrl(a.GetRootDir())

	packages := a.GuessPackages()
	version := a.FindVersion()
	dbs := a.FindDatabases()
	envVars := a.EnvVars()

	services, err := packs.AnalyzeProcfile(a)
	if err != nil {
		fmt.Printf("%s Failed to parse Procfile due to %s\n", common.MsgError, err.Error())
		return nil, err
	}
	err = a.AnalyzeServices(&services)
	if err != nil {
		return nil, err
	}
	packs.RefineServices(&services, envVars, gitBranch, gitURL)

	analysis := &Analysis{
		PackName:           a.GetPack().Name(),
		GitBranch:          gitBranch,
		GitURL:             gitURL,
		ServiceYAMLContext: &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs.Items}},
		DockerfileContext:  &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}},
		Messages:           a.Messages}
	return analysis, nil
}

func (a *Analyzer) AnalyzeServices(services *[]*common.Service) error {
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
