package packs

import (
	"fmt"
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

func (r *Ruby) OutputFolder() string {
	return r.WorkDir
}

func (r *Ruby) Compile() error {
	// we have a ruby app

	// TODO: should check to see if Gemfile has fixed ruby version
	r.Version = "onbuild"

	// port depends on the application server. for now we are going to fix to 3000
	if runsUnicorn, _ := common.GetGemVersion(r.Gemfile, "unicorn", "thin"); runsUnicorn {
		fmt.Println("----> Found non Webrick application server")
		r.Ports = []int{9292}
	} else {
		r.Ports = []int{3000}
	}

	// TODO: add packages based on any other findings in the Gemfile
	r.Packages = common.NewLister()

	// look for DB
	r.Dbs = common.NewLister()
	if hasMysql, _ := common.GetGemVersion(r.Gemfile, "mysql2"); hasMysql {
		fmt.Println("----> Found Mysql")
		r.Dbs.Add("mysql")
	}

	if hasPg, _ := common.GetGemVersion(r.Gemfile, "pg"); hasPg {
		fmt.Println("----> Found PostgreSQL")
		r.Dbs.Add("postgres")
	}

	if hasRedis, _ := common.GetGemVersion(r.Gemfile, "redis"); hasRedis {
		fmt.Println("----> Found Redis")
		r.Dbs.Add("redis")
	}

	return nil
}
