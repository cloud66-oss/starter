package python

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/python/webservers"
)

type Analyzer struct {
	packs.AnalyzerBase
	RequirementsTxt string
	SettingsPy      string
	WSGIFile        string
	PythonPackages  []string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	var hasFound bool
	var err error

	a.RequirementsTxt = a.findRequirementsTxt()
	a.PythonPackages, err = common.PythonPackages(a.RequirementsTxt)
	if err != nil {
		return nil, err
	}

	hasFound, a.WSGIFile = a.findWSGIFile()
	if !hasFound {
		return nil, fmt.Errorf("Could not find WSGI file")
	}

	hasFound, a.SettingsPy = a.findSettingsPy()
	if !hasFound {
		return nil, fmt.Errorf("Could not find settings file")
	}

	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}

	packages := a.GuessPackages()
	version := a.FindVersion()
	dbs, err := a.FindDatabases()
	if err != nil {
		return nil, err
	}
	dbs = a.ConfirmDatabases(dbs)
	envVars := a.EnvVars()

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
		DockerfileContext: &DockerfileContext{
			DockerfileContextBase: packs.DockerfileContextBase{Version: version, Packages: packages},
			RequirementsTxt:       a.RequirementsTxt}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	service := a.GetOrCreateWebService(services)
	var command string
	ports := []*common.PortMapping{common.NewPortMapping()}
	hasFound, server := a.detectWebServer()
	if hasFound {
		// The command was found in the Procfile
		ports[0].Container = server.Port(service.Command)
	} else {
		if common.IsDjangoProject(a.RootDir) {
			command = "python manage.py runserver"
			ports[0].Container = "8000"
		} else {
			//TODO:
		}
	}

	if service.Command == "" {
		service.Command = command
	}

	service.BuildCommand = a.AskForCommand("python manage.py migrate", "build")
	service.DeployCommand = a.AskForCommand("python manage.py migrate", "deployment")

	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	return common.ContainsString(a.PythonPackages, pack)
}

func (a *Analyzer) detectWebServer() (hasFound bool, server packs.WebServer) {
	gunicorn := &webservers.Gunicorn{}
	servers := []packs.WebServer{gunicorn}
	return a.AnalyzerBase.DetectWebServer(a, servers)
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	return packages
}

func (a *Analyzer) FindVersion() string {
	hasFound, version := common.GetPythonVersion()
	return a.ConfirmVersion(hasFound, version, "latest")

}
func (a *Analyzer) FindDatabases() (*common.Lister, error) {
	dbs := common.NewLister()
	settings, err := ioutil.ReadFile(a.SettingsPy)
	if err != nil {
		return nil, err
	}

	dbNames := map[string]string{
		"mysql":                 "mysql",      // Django built-in
		"postgresql_psycopg2":   "postgresql", // Django built-in
		"RedisCache":            "redis",      // https://github.com/niwinz/django-redis or https://github.com/sebleier/django-redis-cache
		"django_mongodb_engine": "mongodb",    // https://github.com/django-nonrel/mongodb-engine
		"mongodb":               "mongodb",    // https://github.com/peterbe/django-mongokit
	}

	databaseDictionaryPattern := regexp.MustCompile(`(?m)^[[:blank:]]*(?:DATABASES|CACHES)[[:blank:]]*=[[:blank:]]*{(?:.|\n)*?^}`)
	databasePattern := regexp.MustCompile(`(?m)^[[:blank:]]*['"](?:ENGINE|BACKEND)['"][[:blank:]]*:[[:blank:]]*['"](?:[[:word:]]*\.)*([[:word:]]*)['"]`)
	for _, dictionary := range databaseDictionaryPattern.FindAllString(string(settings), -1) {
		for _, match := range databasePattern.FindAllStringSubmatch(dictionary, -1) {
			pythonDbName := match[1]
			if pythonDbName == "sqlite3" {
				continue
			}
			db := dbNames[pythonDbName]
			if db == "" {
				return nil, fmt.Errorf("Found not supported database backend %s, aborting", pythonDbName)
			}
			dbs.Add(db)
		}
	}

	return dbs, nil
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{}
}

func (a *Analyzer) findRequirementsTxt() string {
	return common.AskUserWithDefault("Enter path to requirements file", "requirements.txt", a.ShouldPrompt)
}

func (a *Analyzer) findWSGIFile() (bool, string) {
	var found []string
	WSGIPattern := regexp.MustCompile("^(wsgi\\.py|.*\\.wsgi)$")
	filepath.Walk(a.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if WSGIPattern.MatchString(filepath.Base(path)) {
			found = append(found, path)
		}
		return nil
	})

	wsgi := ""
	if len(found) == 1 && common.AskYesOrNo(common.MsgL1, fmt.Sprintf("Found WSGI file %s, confirm?", found[0]), true, a.ShouldPrompt) {
		wsgi = found[0]
	} else if len(found) > 1 && a.ShouldPrompt {
		answer := common.AskMultipleChoices("Found several potential WSGI files. Please choose one:", append(found, "Other"))
		if answer != "Other" {
			wsgi = answer
		}
	}

	if wsgi == "" && a.ShouldPrompt {
		wsgi = common.AskUser("Enter WSGI file path")
	}
	return wsgi != "", wsgi
}

func (a *Analyzer) findSettingsPy() (hasFound bool, path string) {
	hasFound, settingsModule := a.djangoSettingsModule()

	message := "Enter production settings file path"
	if hasFound {
		return true, common.AskUserWithDefault(message, a.module2File(settingsModule), a.ShouldPrompt)
	}
	if a.ShouldPrompt {
		return true, common.AskUser(message + " (e.g 'yourapp/settings.py')")
	}
	return false, ""
}

func (a *Analyzer) djangoSettingsModule() (bool, string) {
	settingsModule := os.Getenv("DJANGO_SETTINGS_MODULE")
	if settingsModule != "" {
		return true, settingsModule
	}

	wsgi, err := ioutil.ReadFile(a.WSGIFile)
	if err != nil {
		return false, ""
	}

	settingsPattern := regexp.MustCompile(`(?m)^[[:blank:]]*os.environ.setdefault\("DJANGO_SETTINGS_MODULE", "(.*)"\)`)
	match := settingsPattern.FindStringSubmatch(string(wsgi))
	if len(match) > 0 {
		return true, match[1]
	}

	return false, ""
}

func (a *Analyzer) module2File(moduleName string) string {
	return strings.Replace(moduleName, ".", string(os.PathSeparator), -1) + ".py"
}
