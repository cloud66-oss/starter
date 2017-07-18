package service_yml

import (
	"path/filepath"

	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/common"
)

type Detector struct {
	packs.PackElement
}

func (d *Detector) Detect(rootDir string) bool {
	return common.FileExists(filepath.Join(rootDir, "service.yml"))
}


