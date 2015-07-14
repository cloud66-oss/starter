package node

import "github.com/cloud66/starter/common"

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
