package node

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Detector struct {
	packs.Detector
}

func (d *Detector) PackName() string {
	return "node"
}

func (d *Detector) Detect(rootDir string) bool {
	return common.FileExists(filepath.Join(rootDir, "package.json"))
}
