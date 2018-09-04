package packs

import "github.com/cloud66-oss/starter/common"

type DockerfileContextBase struct {
	Version          string
	Framework        string
	Packages         *common.Lister
	FrameworkVersion string
}
