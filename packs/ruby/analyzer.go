package ruby

import (
	"fmt"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
)

type Analyzer struct {
	packs.AnalyzerBase
	Gemfile string
}

func (a *Analyzer) Name() string {
	return "ruby"
}

func (a *Analyzer) Analyze() error {
	messages := common.NewLister()

	foundRuby, rubyVersion := common.GetRubyVersion(a.Gemfile)
	if foundRuby {
		a.Version = fmt.Sprintf("%s-onbuild", rubyVersion)
	} else {
		rubyVersion = common.AskUser("Can't find Ruby version from Gemfile:", "default")
		if rubyVersion == "default" {
			a.Version = a.defaultVersion()
		} else {
			a.Version = fmt.Sprintf("%s-onbuild", rubyVersion)
		}
	}

	service := &common.Service{Name: "web"}

	isRails, _ := common.GetGemVersion(a.Gemfile, "rails")

	// port depends on the application server. for now we are going to fix to 3000
	if runsUnicorn, _ := common.GetGemVersion(a.Gemfile, "unicorn", "thin"); runsUnicorn {
		fmt.Println(common.MsgL2, "----> Found non Webrick application server", common.MsgReset)
		// The command here will be found in the Procfile
		service.Ports = []string{"9292:80:443"}
	} else {
		if isRails {
			service.Command = "bundle exec rails s _env:RAILS_ENV"
			service.Ports = []string{"3000:80:443"}
		} else {
			service.Command = "bundle exec rackup s _env:RACK_ENV"
			service.Ports = []string{"9292:80:443"}
		}
	}

	// add packages based on any other findings in the Gemfile
	a.Packages = common.NewLister()
	if hasRmagick, _ := common.GetGemVersion(a.Gemfile, "rmagick"); hasRmagick {
		fmt.Println(common.MsgL2, "----> Found Image Magick", common.MsgReset)
		a.Packages.Add("imagemagick", "libmagickwand-dev")
	}

	if hasSqlite, _ := common.GetGemVersion(a.Gemfile, "sqlite"); hasSqlite {
		a.Packages.Add("libsqlite3-dev")
		fmt.Println(common.MsgL2, "----> Found sqlite", common.MsgReset)
	}

	// look for DB
	dbs := common.NewLister()
	if hasMysql, _ := common.GetGemVersion(a.Gemfile, "mysql2"); hasMysql {
		fmt.Println(common.MsgL2, "----> Found Mysql", common.MsgReset)
		dbs.Add("mysql")
	}

	if hasPg, _ := common.GetGemVersion(a.Gemfile, "pg"); hasPg {
		fmt.Println(common.MsgL2, "----> Found PostgreSQL", common.MsgReset)
		dbs.Add("postgresql")
	}

	if hasRedis, _ := common.GetGemVersion(a.Gemfile, "redis"); hasRedis {
		fmt.Println(common.MsgL2, "----> Found Redis", common.MsgReset)
		dbs.Add("redis")
	}

	if hasMongoDB, _ := common.GetGemVersion(a.Gemfile, "mongo", "mongo_mapper", "dm-mongo-adapter", "mongoid"); hasMongoDB {
		fmt.Println(common.MsgL2, "----> Found MongoDB", common.MsgReset)
		dbs.Add("mongodb")
	}

	if hasElasticsearch, _ := common.GetGemVersion(a.Gemfile, "elasticsearch", "tire", "flex", "chewy"); hasElasticsearch {
		fmt.Println(common.MsgL2, "----> Found Elasticsearch", common.MsgReset)
		dbs.Add("elasticsearch")
	}

	if hasDatabaseYaml := common.FileExists("config/database.yml"); hasDatabaseYaml {
		fmt.Println(common.MsgL2, "----> Found config/database.yml", common.MsgReset)
		messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"database.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}

	if hasMongoIdYaml := common.FileExists("config/mongoid.yml"); hasMongoIdYaml {
		fmt.Println(common.MsgL2, "----> Found config/mongoid.yml", common.MsgReset)
		messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"mongoid.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}

	a.Context = &common.ParseContext{
		Services: []*common.Service{service},
		Dbs:      dbs.Items,
		EnvVars: []*common.EnvVar{
			&common.EnvVar{Key: "RAILS_ENV", Value: a.Environment},
			&common.EnvVar{Key: "RACK_ENV", Value: a.Environment}},
		Messages: messages.Items}

	service.EnvVars = a.Context.EnvVars

	return nil
}

func (a *Analyzer) defaultVersion() string {
	return "onbuild"
}
