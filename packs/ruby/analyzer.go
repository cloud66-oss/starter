package ruby

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/ruby/webservers"
)

type Analyzer struct {
	packs.AnalyzerBase
	Gemfile string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.Gemfile = filepath.Join(a.RootDir, "Gemfile")
	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}

	version := a.FindVersion()
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()
	packages := a.GuessPackages()

	services, err := a.AnalyzeServices(a, envVars, gitBranch, gitURL, buildRoot)
	if err != nil {
		return nil, err
	}

	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:  a.GetPack().Name(),
			GitBranch: gitBranch,
			GitURL:    gitURL,
			Messages:  a.Messages},
		ServiceYAMLContext: &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs.Items}},
		DockerfileContext:  &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	service := a.GetOrCreateWebService(services)
	service.Ports = []*common.PortMapping{common.NewPortMapping()}
	isRails, _ := common.GetGemVersion(a.Gemfile, "rails")

	if service.Command == "" {
		if isRails {
			service.Command = "bundle exec rails server -e _env:RAILS_ENV"
			service.Ports[0].Container = "3000"
		} else {
			service.Command = "bundle exec rackup -E _env:RACK_ENV"
			service.Ports[0].Container = "9292"
		}
	} else {
		var err error
		hasFoundServer, server := a.detectWebServer(service.Command)
		service.Ports[0].Container, err = a.FindPort(hasFoundServer, server, &service.Command)

		if err != nil {
			return err
		}
	}

	if isRails {
		service.BuildCommand = a.AskForCommand("bundle exec rake db:schema:load", "build")
		service.DeployCommand = a.AskForCommand("bundle exec rake db:migrate", "deployment")
	} else {
		service.BuildCommand = a.AskForCommand("", "build")
		service.DeployCommand = a.AskForCommand("", "deployment")
	}

	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	hasFound, _ := common.GetGemVersion(a.Gemfile, pack)
	return hasFound
}

func (a *Analyzer) detectWebServer(command string) (hasFound bool, server packs.WebServer) {
	unicorn := &webservers.Unicorn{}
	thin := &webservers.Thin{}
	servers := []packs.WebServer{unicorn, thin}
	return a.AnalyzerBase.DetectWebServer(a, command, servers)
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	if hasRmagick, _ := common.GetGemVersion(a.Gemfile, "rmagick", "refile-mini_magick", "mini_magick"); hasRmagick {
		fmt.Println(common.MsgL2, "----> Found Image Magick", common.MsgReset)
		packages.Add("imagemagick", "libmagickwand-dev")
	}

	if hasSqlite, _ := common.GetGemVersion(a.Gemfile, "sqlite"); hasSqlite {
		packages.Add("libsqlite3-dev")
		fmt.Println(common.MsgL2, "----> Found sqlite", common.MsgReset)
	}

	if hasMemcache, _ := common.GetGemVersion(a.Gemfile, "dalli"); hasMemcache {
		packages.Add("memcached")
		fmt.Println(common.MsgL2, "----> Found Memcache", common.MsgReset)
	}
	return packages
}

func (a *Analyzer) FindVersion() string {
	foundRuby, rubyVersion := common.GetRubyVersion(a.Gemfile)
	return a.ConfirmVersion(foundRuby, rubyVersion, "latest")
}

func (a *Analyzer) FindDatabases() *common.Lister {
	dbs := common.NewLister()
	if hasMysql, _ := common.GetGemVersion(a.Gemfile, "mysql2"); hasMysql {
		dbs.Add("mysql")
	}

	if hasPg, _ := common.GetGemVersion(a.Gemfile, "pg"); hasPg {
		dbs.Add("postgresql")
	}

	if hasRedis, _ := common.GetGemVersion(a.Gemfile, "redis", "redis-rails"); hasRedis {
		dbs.Add("redis")
	}

	if hasMongoDB, _ := common.GetGemVersion(a.Gemfile, "mongo", "mongo_mapper", "dm-mongo-adapter", "mongoid"); hasMongoDB {
		dbs.Add("mongodb")
	}

	if hasElasticsearch, _ := common.GetGemVersion(a.Gemfile, "elasticsearch", "tire", "flex", "chewy"); hasElasticsearch {
		dbs.Add("elasticsearch")
	}

	if hasDatabaseYaml := common.FileExists("config/database.yml"); hasDatabaseYaml {
		fmt.Println(common.MsgL2, "----> Found config/database.yml", common.MsgReset)
		a.Messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"database.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}

	if hasMongoIdYaml := common.FileExists("config/mongoid.yml"); hasMongoIdYaml {
		fmt.Println(common.MsgL2, "----> Found config/mongoid.yml", common.MsgReset)
		a.Messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"mongoid.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{
		&common.EnvVar{Key: "RAILS_ENV", Value: a.Environment},
		&common.EnvVar{Key: "RACK_ENV", Value: a.Environment}}
}
