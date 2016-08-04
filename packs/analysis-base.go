package packs

import "github.com/cloud66/starter/common"

type AnalysisBase struct {
	PackName string

	GitURL    string
	GitBranch string

	Messages common.Lister
}
