package packs

import "github.com/cloud66/starter/common"

type DockerfileContextBase struct {
	Version  string
	Framework string
	Packages *common.Lister
}
