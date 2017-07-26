package service_yml_to_kubes

import (
	"github.com/cloud66/starter/packs"
	"path/filepath"
)

type Analyzer struct{
	packs.AnalyzerBase
	ServiceYml string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.ServiceYml = filepath.Join(a.RootDir, "service.yml")
	gitURL, gitBranch, _, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}
	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:  a.GetPack().Name(),
			GitBranch: gitBranch,
			GitURL:    gitURL,
			Messages:  a.Messages},
	}
	return analysis, nil
}
