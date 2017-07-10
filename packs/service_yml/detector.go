package service_yml

import (
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/common"
	"path/filepath"
)

type Detector struct {
	packs.PackElement
}

func (d *Detector) Detect(rootDir string) bool {
	return common.FileExists(filepath.Join(rootDir, "service.yml"))
}
