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
	foundNode, nodeVersion := common.GetNodeVersion(a.PackageJSON)

	if foundNode {
		a.Version = fmt.Sprintf("%s-onbuild", nodeVersion)
	} else {
		//TODO: I'm not sure this will happen!
		nodeVersion = common.AskUser("Can't find Node version from package.json:", "default")
		if nodeVersion == "default" {
			a.Version = a.defaultVersion()
		} else {
			a.Version = fmt.Sprintf("%s-onbuild", nodeVersion)
		}
	}

	messages := common.NewLister()

	service := &common.Service{Name: "web"}

	if runsExpress, _ := common.GetDependencyVersion(a.PackageJSON, "express"); runsExpress {
		fmt.Println(common.MsgL2, "----> Found Express", common.MsgReset)
	}

	if hasScript, script := common.GetScriptsStart(a.PackageJSON); hasScript {
		fmt.Println(common.MsgL2, "----> Found Script:", script, common.MsgReset)
	}

	dbs := common.NewLister()

	a.Context = &common.ParseContext{
		Services: []*common.Service{service},
		Dbs:      dbs.Items,
		EnvVars:  []*common.EnvVar{},
		Messages: messages.Items}

	service.EnvVars = a.Context.EnvVars

	return nil
}

func (a *Analyzer) defaultVersion() string {
	return "onbuild"
}
