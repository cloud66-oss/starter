package ruby

import "github.com/cloud66-oss/starter/packs"

type Analysis struct {
	packs.AnalysisBase
	DockerComposeYAMLContext *DockerComposeYAMLContext
	ServiceYAMLContext       *ServiceYAMLContext
	DockerfileContext        *DockerfileContext
}
