package packs

import "github.com/cloud66/starter/common"

type Analyzer interface {
	FillServices(*[]*common.Service) error
	HasPackage(pack string) bool
	GuessFramework() string
}
