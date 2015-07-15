package node

import (
	"fmt"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Analyzer struct {
	packs.AnalyzerBase
	PackageJSON string
}

func (a *Analyzer) Name() string {
	return "node"
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
		//TODO: I'm not sure this will happen!
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
