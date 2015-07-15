package ruby

import "github.com/cloud66/starter/packs"

type Pack struct {
}

func (p *Pack) Name() string {
	return "ruby"
}

func (p *Pack) Detector() packs.Detector {
	return &Detector{PackElement: packs.PackElement{Pack: p}}
}

func (p *Pack) Analyzer(rootDir string, environment string) packs.Analyzer {
	return &Analyzer{
		AnalyzerBase: packs.AnalyzerBase{PackElement: packs.PackElement{Pack: p},
			RootDir:     rootDir,
			Environment: environment}}
}
