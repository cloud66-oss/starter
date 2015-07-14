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

func (d *Detector) Detect(rootDir string) bool {
	d.packageJSON = filepath.Join(rootDir, "package.json")
	return common.FileExists(d.packageJSON)
}

func (d *Detector) Analyzer(rootDir string, environment string) packs.Analyzer {
	return &Analyzer{PackageJSON: d.packageJSON, AnalyzerBase: packs.AnalyzerBase{RootDir: rootDir, Environment: environment}}
}
