package packs

import "github.com/cloud66/starter/common"

type DockerComposeYAMLContextBase struct {
	Services []*common.Service
	Dbs      []*common.Database
}
