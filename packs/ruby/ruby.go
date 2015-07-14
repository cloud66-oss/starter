package ruby

import "github.com/cloud66/starter/common"

type Ruby struct {
	WorkDir     string
	Environment string

	Gemfile  string
	Version  string
	Packages *common.Lister
}

func (r *Ruby) Name() string {
	return "ruby"
}

func (r *Ruby) PackVersion() string {
	return "0.1"
}

func (r *Ruby) OutputFolder() string {
	return r.WorkDir
}
