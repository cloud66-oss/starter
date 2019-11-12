package node

import (
	"path/filepath"

	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/packs"
)

type Detector struct {
	packs.PackElement
}

func (d *Detector) Detect(rootDir string) bool {
	if common.FileExists(filepath.Join(rootDir, "Gemfile")) || common.FileExists(filepath.Join(rootDir, "config", "database.yml")) {
		return false
	}
	return common.FileExists(filepath.Join(rootDir, "package.json"))
}
