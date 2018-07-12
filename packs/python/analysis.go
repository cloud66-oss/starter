package python

import "github.com/cloud66-oss/starter/packs"

type Analysis struct {
	packs.AnalysisBase

	ServiceYAMLContext *ServiceYAMLContext
	DockerfileContext  *DockerfileContext
}
