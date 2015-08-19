package python

import "github.com/cloud66/starter/packs"

type Analysis struct {
	packs.AnalysisBase

	ServiceYAMLContext *ServiceYAMLContext
	DockerfileContext  *DockerfileContext
}
