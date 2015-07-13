package packs

import (
	"github.com/cloud66/starter/common"
)

type Pack interface {
	Name() string
	PackVersion() string
	Compile() (*common.ParseContext, error)
	OutputFolder() string
	DefaultVersion() string
}
