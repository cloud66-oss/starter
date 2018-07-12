package packs

import "github.com/cloud66-oss/starter/common"

type Analyzer interface {
	FillServices(*[]*common.Service) error
	HasPackage(pack string) bool
	GuessFramework() string
}
