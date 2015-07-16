package packs

import "github.com/cloud66/starter/common"

type Pack interface {
	Name() string
	Detector() Detector
	Analyze(rootDir string, environment string) error
	WriteDockerfile(string, string, bool) error
	WriteServiceYAML(string, string, bool) error
	GetMessages() []string
}

type PackBase struct {
	Messages *common.Lister
}

type PackElement struct {
	Pack Pack
}

func (e *PackElement) GetPack() Pack {
	return e.Pack
}
