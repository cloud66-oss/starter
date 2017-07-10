package docker_compose

import (
	"github.com/cloud66/starter/packs"
	"path/filepath"
)

type Analyzer struct{
	packs.AnalyzerBase
	DockerCompose string
}


func (a *Analyzer) Analyze() (*Analysis, error) {
	a.DockerCompose = filepath.Join(a.RootDir, "docker-compose.yml")
	gitURL, gitBranch, _, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}


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