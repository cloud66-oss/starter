package node

import (
	"fmt"
	"path/filepath"
	"strings"
	"strconv"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/blang/semver"
)

type Analyzer struct {
	packs.AnalyzerBase
	PackageJSON string
	MeteorRelease string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.PackageJSON = filepath.Join(a.RootDir, "package.json")
	a.MeteorRelease = filepath.Join(a.RootDir, ".meteor/release")



	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()
	packages := a.GuessPackages()
	framework := a.GuessFramework()

	//TODO make this more generic
	if framework == "keystone" {
		dbs = append(dbs, common.Database{Name: "mongodb", DockerImage: "mongo"})
	}

	if framework == "meteor" {
		dbs = append(dbs, common.Database{Name: "mongodb", DockerImage: "mongo"})
	}


	framework_version := a.GuessFrameworkVersion()
	supported_versions := a.FindVersion()
	version := supported_versions[len(supported_versions)-1]
	//a.CheckNotSupportedPackages(packages)

	services, err := a.AnalyzeServices(a, envVars, gitBranch, gitURL, buildRoot)

	// inject all the services with the databases used in the infrastructure
	listOfStartCommands := []string {}

	for _, service := range services {
		listOfStartCommands = append(listOfStartCommands, service.Command)
		service.Databases = dbs
	}

	if err != nil {
		return nil, err
	}


	listOfDatabases := []string {}

	for _, database := range dbs {
		listOfDatabases = append(listOfDatabases, database.Name)
	}

	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:         a.GetPack().Name(),
			GitBranch:        gitBranch,
			GitURL:           gitURL,
			Framework:        framework,
			FrameworkVersion: framework_version,
			LanguageVersion:  version,
			SupportedLanguageVersions: supported_versions,
			Databases:			listOfDatabases,
			ListOfStartCommands:	listOfStartCommands,
			Messages:         a.Messages},
		DockerComposeYAMLContext: &DockerComposeYAMLContext{packs.DockerComposeYAMLContextBase{Services: services, Dbs: dbs}},
		ServiceYAMLContext:       &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs}},
		DockerfileContext:        &DockerfileContext{packs.DockerfileContextBase{Version: version, Framework: framework, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	//has a procfile
	if len(*services) > 0 {

		for _, service := range *services {
			port := 3001
			if service.Name == "web" || service.Name == "custom_web" {
				service.Ports = []*common.PortMapping{common.NewPortMapping()}
				service.Ports[0].Container = "3000"
				service.EnvVars = []*common.EnvVar{common.NewEnvMapping("PORT", "3000")}
			} else {
				service.Ports = []*common.PortMapping{common.NewInternalPortMapping(strconv.Itoa(port))}
				service.EnvVars = []*common.EnvVar{common.NewEnvMapping("PORT", strconv.Itoa(port))}
				port = port +1
			}
		}
	}

	//no procfile
	if len(*services) == 0 {
		var service *common.Service
		service = &common.Service{Name: "web"}
		service.Ports = []*common.PortMapping{common.NewPortMapping()}
		service.Command = "npm start"
		if hasScript, script := common.GetScriptsStart(a.PackageJSON); hasScript {

			common.PrintlnL2("Found Script: %s", script)
			service.Command = script
		}
		if a.GuessFramework() == "meteor" {
			service.Command = "node main.js"
			common.PrintlnL2("Change Script for Meteor: %s", service.Command)
		}

		service.Ports[0].Container = "3000"
		service.EnvVars = []*common.EnvVar{common.NewEnvMapping("PORT", "3000")}
		*services = append(*services, service)
	}

	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	hasFound, _ := common.GetDependencyVersion(a.PackageJSON, pack)
	return hasFound
}

func (a *Analyzer) GetPackageVersion(pack string) string {
	hasFound, version := common.GetDependencyVersion(a.PackageJSON, pack)
	if hasFound {
		v1, err := semver.Make(strings.Replace(strings.Trim(version, "^>=~"), "x", "0", -1))
		if err != nil {
		  return ""
		}
		version = fmt.Sprintf("%d.%d.%d", v1.Major, v1.Minor, v1.Patch)
		return version
	} else {
		return ""
	}
}

func (a *Analyzer) GuessFramework() string {
	for _, framework := range common.GetSupportedNodeFrameworks() {
		if a.HasPackage(framework) {
			if framework == "meteor-node-stubs" {
				return "meteor"		
			}
			return framework
		}
	}

	return ""
}

func (a *Analyzer) GetMeteorVersion() string {
	_, meteorVersion := common.GetMeteorVersion(a.MeteorRelease)
	return meteorVersion
}



func (a *Analyzer) GuessFrameworkVersion() string {
	for _, framework := range common.GetSupportedNodeFrameworks() {
		if a.HasPackage(framework) {
			if framework == "meteor-node-stubs" {
				return a.GetMeteorVersion()		
			}
			return a.GetPackageVersion(framework)
		}
	}
	return ""
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	return packages
}

func (a *Analyzer) FindVersion() []string {
	_, nodeVersions := common.GetNodeVersion(a.PackageJSON)
	//return a.ConfirmVersion(foundNode, nodeVersions[0], common.GetDefaultNodeVersion())
	return nodeVersions
}

func (a *Analyzer) FindDatabases() []common.Database {
	dbs := []common.Database{}
	if hasMysql, _ := common.GetNodeDatabase(a.PackageJSON, "mysql", "loopback-connector-mysql"); hasMysql {
		dbs = append(dbs, common.Database{Name: "mysql", DockerImage: "mysql"})
	}
	if hasMongo, _ := common.GetNodeDatabase(a.PackageJSON, "mongoose", "mongodb", "loopback-connector-mongodb"); hasMongo {
		dbs = append(dbs, common.Database{Name: "mongodb", DockerImage: "mongo"})
	}
	if hasPostgres, _ := common.GetNodeDatabase(a.PackageJSON, "pg", "loopback-connector-postgresql"); hasPostgres {
		dbs = append(dbs, common.Database{Name: "postgresql", DockerImage: "postgres"})
	}
	if hasRedis, _ := common.GetNodeDatabase(a.PackageJSON, "redis", "ioredis", "loopback-connector-redis", "loopback-connector-kv-redis"); hasRedis {
		dbs = append(dbs, common.Database{Name: "redis", DockerImage: "redis"})
	}
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}
