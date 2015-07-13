package node

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Detector struct {
	packs.Detector
	packageJSON string
}

func (d *Detector) Name() string {
	return "Node"
}

func (d *Detector) Detect(root string) bool {
	d.packageJSON = filepath.Join(root, "package.json")
	return common.FileExists(d.packageJSON)
}

func (d *Detector) Analyzer(root string, environment string) packs.Pack {
	return &Node{PackageJson: d.packageJSON, WorkDir: root, Environment: environment}
}
