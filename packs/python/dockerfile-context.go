package python

import "github.com/cloud66-oss/starter/packs"

type DockerfileContext struct {
	packs.DockerfileContextBase
	RequirementsTxt string
}
