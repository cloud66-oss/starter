package packs

import "github.com/cloud66/starter/common"

type AnalysisBase struct {
	PackName string
	LanguageVersion string
	GitURL    string
	GitBranch string
	Framework string
	FrameworkVersion string
	Messages common.Lister
}
