package ruby

import "github.com/cloud66/starter/packs"

type Analysis struct {
	packs.AnalysisBase
	DockerComposeYAMLContext *DockerComposeYAMLContext
	ServiceYAMLContext *ServiceYAMLContext
	DockerfileContext  *DockerfileContext

}
