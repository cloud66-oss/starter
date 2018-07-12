package ruby

import (
	"path/filepath"

	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/packs"
)

type Detector struct {
	packs.PackElement
}

func (d *Detector) Detect(rootDir string) bool {
	return common.FileExists(filepath.Join(rootDir, "Gemfile"))
}
