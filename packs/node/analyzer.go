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

func (a *Analyzer) Analyze() error {
	service := &common.Service{Name: "web"}
	if runsExpress, _ := common.GetDependencyVersion(a.PackageJSON, "express"); runsExpress {
		fmt.Println(common.MsgL2, "----> Found Express", common.MsgReset)
	}
	if hasScript, script := common.GetScriptsStart(a.PackageJSON); hasScript {
		fmt.Println(common.MsgL2, "----> Found Script:", script, common.MsgReset)
	}

	a.Packages = a.GuessPackages()
	a.Version = a.FindVersion()
	a.Context = &common.ParseContext{
		Services: []*common.Service{service},
		Dbs:      a.AnalyzeDatabases().Items,
		EnvVars:  []*common.EnvVar{}}

	return nil
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
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

func (a *Analyzer) defaultVersion() string {
	return "onbuild"
}

func (a *Analyzer) AnalyzeDatabases() *common.Lister {
	dbs := common.NewLister()
	return dbs
}
