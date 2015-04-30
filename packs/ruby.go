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

func (r *Ruby) Compile() (*common.ParseContext, error) {
	// we have a ruby app

	// TODO: should check to see if Gemfile has fixed ruby version
	r.Version = "onbuild"

	service := &common.Service{Name: "web"}

	// port depends on the application server. for now we are going to fix to 3000
	if runsUnicorn, _ := common.GetGemVersion(r.Gemfile, "unicorn", "thin"); runsUnicorn {
		fmt.Println("----> Found non Webrick application server")
		service.Ports = []int{9292}
	} else {
		service.Ports = []int{3000}
	}

	// TODO: read and parse Procfile and use that to build the services

	// add packages based on any other findings in the Gemfile
	r.Packages = common.NewLister()
	if hasRmagick, _ := common.GetGemVersion(r.Gemfile, "rmagick"); hasRmagick {
		fmt.Println("----> Found Image Magick")
		r.Packages.Add("imagemagick", "libmagickwand-dev")
	}

	if hasSqlite, _ := common.GetGemVersion(r.Gemfile, "sqlite"); hasSqlite {
		fmt.Println("----> Found sqlite")
		r.Packages.Add("libsqlite3-dev")
	}

	// look for DB
	dbs := common.NewLister()
	if hasMysql, _ := common.GetGemVersion(r.Gemfile, "mysql2"); hasMysql {
		fmt.Println("----> Found Mysql")
		dbs.Add("mysql")
	}

	if hasPg, _ := common.GetGemVersion(r.Gemfile, "pg"); hasPg {
		fmt.Println("----> Found PostgreSQL")
		dbs.Add("postgres")
	}

	if hasRedis, _ := common.GetGemVersion(r.Gemfile, "redis"); hasRedis {
		fmt.Println("----> Found Redis")
		dbs.Add("redis")
	}

	parseContext := &common.ParseContext{Services: []*common.Service{service}, Dbs: dbs}

	return parseContext, nil
}
