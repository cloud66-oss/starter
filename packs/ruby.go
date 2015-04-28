package packs

import (
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type Ruby struct {
	WorkDir string

	Gemfile  string
	Version  string
	Packages *common.Lister
	Ports    []int
	Dbs      *common.Lister

	GitRepo   string
	GitBranch string
	Command   string
}

func (r *Ruby) Name() string {
	return "ruby"
}

func (r *Ruby) PackVersion() string {
	return "0.1"
}

func (r *Ruby) Detect() (bool, error) {
	r.Gemfile = filepath.Join(r.WorkDir, "Gemfile")

	// TODO: fetch git url and branch from the director
	return common.FileExists(r.Gemfile), nil
}

func (r *Ruby) Compile() error {
	// we have a ruby app

	// TODO: should check to see if Gemfile has fixed ruby version
	r.Version = "onbuild"

	// TODO: port depends on the application server. for now we are going to fix to 3000
	r.Ports = []int{3000}

	// TODO: add packages based on any other findings in the Gemfile
	r.Packages = common.NewLister()

	// TODO: look for DB
	r.Dbs = common.NewLister()

	return nil
}
