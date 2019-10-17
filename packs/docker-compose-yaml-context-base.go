package packs

import "github.com/cloud66-oss/starter/common"

type DockerComposeYAMLContextBase struct {
	Services []*common.Service
	Dbs      []common.Database
}
