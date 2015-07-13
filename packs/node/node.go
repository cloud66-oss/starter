package node

import (
	"fmt"

	"github.com/cloud66/starter/common"
)

type Node struct {
	WorkDir     string
	Environment string

	PackageJson string
	Version     string
	Packages    *common.Lister
}

type PackageJson interface{}

func (r *Node) Name() string {
	return "node"
}

func (r *Node) PackVersion() string {
	return "0.1"
}

func (r *Node) OutputFolder() string {
	return r.WorkDir
}

func (r *Node) DefaultVersion() string {
	return "onbuild"
}

func (r *Node) Compile() (*common.ParseContext, error) {
	// we have a Node app
	foundNode, nodeVersion := common.GetNodeVersion(r.PackageJson)

	if foundNode {
		r.Version = fmt.Sprintf("%s-onbuild", nodeVersion)
	} else {
		//TODO: I'm not sure this will happen!
		nodeVersion = common.AskUser("Can't find Node version from package.json:", "default")
		if nodeVersion == "default" {
			r.Version = r.DefaultVersion()
		} else {
			r.Version = fmt.Sprintf("%s-onbuild", nodeVersion)
		}
	}

	messages := common.NewLister()

	service := &common.Service{Name: "web"}

	if runsExpress, _ := common.GetDependencyVersion(r.PackageJson, "express"); runsExpress {
		fmt.Println(common.MsgL2, "----> Found Express", common.MsgReset)
	}

	if hasScript, script := common.GetScriptsStart(r.PackageJson); hasScript {
		fmt.Println(common.MsgL2, "----> Found Script:", script, common.MsgReset)
	}

	// look for DB
	dbs := common.NewLister()
	// from package.json dependencies

	parseContext := &common.ParseContext{
		Services: []*common.Service{service},
		Dbs:      dbs.Items,
		EnvVars: []*common.EnvVar{
			&common.EnvVar{Key: "RAILS_ENV", Value: r.Environment},
			&common.EnvVar{Key: "RACK_ENV", Value: r.Environment}},
		Messages: messages.Items}

	service.EnvVars = parseContext.EnvVars

	return parseContext, nil
}
