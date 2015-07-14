package ruby

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Detector struct {
	packs.Detector
	gemfile string
}

func (d *Detector) Name() string {
	return "Ruby"
}

func (d *Detector) Detect(rootDir string) bool {
	d.gemfile = filepath.Join(rootDir, "Gemfile")
	return common.FileExists(d.gemfile)
}

func (d *Detector) Analyzer(rootDir string, environment string) packs.Analyzer {
	return &Analyzer{Gemfile: d.gemfile, AnalyzerBase: packs.AnalyzerBase{RootDir: rootDir, Environment: environment}}
}
