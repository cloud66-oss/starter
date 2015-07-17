package packs

import "github.com/cloud66/starter/common"

type ServiceYAMLContextBase struct {
	Services []*common.Service
	Dbs      []string
}
