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

func (d *Detector) Detect(root string) bool {
	d.gemfile = filepath.Join(root, "Gemfile")
	return common.FileExists(d.gemfile)
}

func (d *Detector) Analyzer(root string, environment string) packs.Pack {
	return &Ruby{Gemfile: d.gemfile, WorkDir: root, Environment: environment}
}
