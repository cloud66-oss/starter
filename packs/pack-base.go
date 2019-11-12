package packs

import "github.com/cloud66-oss/starter/common"

const (
	btrSuffix                  = "generic"
	genericStencilTemplatePath = "https://raw.githubusercontent.com/cloud66/stencils/{{.branch}}/"
	genericGithubURL           = "https://github.com/cloud66/stencils.git"
)

type PackBase struct {
	Messages *common.Lister
}

func GenericBundleSuffix() string {
	return btrSuffix
}

func GenericTemplateRepository() string {
	return genericStencilTemplatePath
}

func GithubURL() string {
	return genericGithubURL
}
