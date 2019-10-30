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
	packageJSONFileLocation := filepath.Join(rootDir, "package.json")
	if !common.FileExists(packageJSONFileLocation) {
		return false
	}

	hasFound, _ := common.GetDependencyVersion(packageJSONFileLocation, common.GetSupportedNodeFrameworks()...)
	return hasFound
}
