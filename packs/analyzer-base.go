package packs

import (
	"bufio"
	"fmt"
	"github.com/cloud66-oss/starter/common"
	"os"
	"path/filepath"
	"strings"
)

type AnalyzerBase struct {
	PackElement

	RootDir      string
	Environment  string
	ShouldPrompt bool

	GitURL    string
	GitBranch string

	Messages common.Lister
}

func (a *AnalyzerBase) ProjectMetadata() (string, string, string, error) {
	hasGit := common.HasGit(a.RootDir)
	if hasGit {
		gitURL := common.RemoteGitUrl(a.RootDir)
		gitBranch := common.LocalGitBranch(a.RootDir)
		buildRoot, err := common.PathRelativeToGitRoot(a.RootDir)
		if err != nil {
			return "", "", "", err
		} else {
			return gitURL, gitBranch, buildRoot, nil
		}
	}
	if a.GitURL != "" && a.GitBranch != "" {
		return a.GitURL, a.GitBranch, ".", nil
	} else {
		return "", "", ".", nil
	}
}

func (a *AnalyzerBase) ConfirmDatabases(foundDbs []common.Database) []common.Database {
	var dbs []common.Database
	message := "Found these databases: "
	for i, db := range foundDbs {
		if i != 0 {
			message = message + ", "
		}
		message = message + db.Name
	}

	var dbnames []string
	if !a.ShouldPrompt {
		return foundDbs
	}
	message = message + ". \nList the databases that you desire in your application as a comma separated list."
	common.PrintlnL1(message)
	common.PrintlnL1("  Use 'found' to select the found ones")
	common.PrintlnL1("  Example: 'mysql, redis' ")
	common.PrintL1("> ")
	reader := bufio.NewReader(os.Stdin)
	wantedDbs, err := reader.ReadString('\n')

	if err == nil {
		listOtherDbs := strings.Split(wantedDbs, ",")
		for _, db := range listOtherDbs {
			db = strings.TrimSpace(db)
			if db == "" {
				continue
			}
			dbnames = append(dbnames, db)
		}
		if len(dbnames) == 0 {
			return dbs
		} else {
			if len(dbnames) == 1 && strings.ToLower(dbnames[0]) == "found" {
				return foundDbs
			} else {
				for _, newDb := range dbnames {
					dbs = append(dbs, common.Database{Name: newDb, DockerImage: newDb})
				}
			}
		}
	}

	return dbs
}

func (a *AnalyzerBase) ConfirmVersion(found bool, version string, defaultVersion string) string {
	message := fmt.Sprintf("Found %s version %s, confirm?", a.GetPack().Name(), version)
	if found && common.AskYesOrNo(message, true, a.ShouldPrompt) {
		return version
	}
	return common.AskUserWithDefault(fmt.Sprintf("Enter %s version:", a.GetPack().Name()), defaultVersion, a.ShouldPrompt)
}

func (a *AnalyzerBase) CheckNotSupportedPackages(packages *common.Lister) {
	if packages.Contains("memcached") {
		a.Messages.Add("Memcached was detected but is not currently supported. Please use Redis instead.")
	}
}

func (b *AnalyzerBase) AnalyzeServices(a Analyzer, envVars []*common.EnvVar, gitBranch string, gitURL string, buildRoot string) ([]*common.Service, error) {
	services, err := b.analyzeProcfile()

	if err != nil {
		common.PrintlnError("Failed to parse Procfile due to %s", err.Error())
		return nil, err
	}

	err = a.FillServices(&services)
	if err != nil {
		return nil, err
	}

	b.refineServices(&services)
	b.inheritProjectContext(&services, envVars, gitBranch, gitURL, buildRoot)
	return services, nil
}

func (b *AnalyzerBase) DetectWebServer(a Analyzer, command string, servers []WebServer) (hasFound bool, webserver WebServer) {
	for _, server := range servers {
		for _, name := range server.Names() {
			if a.HasPackage(name) || strings.HasPrefix(command, name) {
				common.PrintlnL2("Found %s", name)
				return true, server
			}
		}
	}
	return false, nil
}

func (a *AnalyzerBase) FindPort(hasFoundServer bool, server WebServer, command *string) (string, error) {
	if hasFoundServer {
		return server.Port(command), nil
	}

	withoutPortEnvVar := common.RemovePortIfEnvVar(*command)
	hasFound, port := common.ParsePort(withoutPortEnvVar)
	if hasFound {
		*command = withoutPortEnvVar
		return port, nil
	}

	if !a.ShouldPrompt {
		return "", fmt.Errorf("Could not find port to open corresponding to command '%s'", *command)
	}
	return common.AskUser(fmt.Sprintf("Which port to open to run web service with command '%s'?", *command)), nil
}

func (a *AnalyzerBase) analyzeProcfile() ([]*common.Service, error) {
	services := []*common.Service{}
	procfilePath := filepath.Join(a.RootDir, "Procfile")
	if !common.FileExists(procfilePath) {
		a.Messages.Add("No Procfile was detected. It is strongly advised to add one in order to specify the commands to run your services.")
		return services, nil
	}

	common.PrintlnL2("Parsing Procfile")
	procs, err := common.ParseProcfile(procfilePath)
	if err != nil {
		return nil, err
	}

	for _, proc := range procs {
		common.PrintlnL2("Found Procfile item %s", proc.Name)
		services = append(services, &common.Service{Name: proc.Name, Command: proc.Command})
	}
	return services, nil
}

func (a *AnalyzerBase) GetOrCreateWebService(services *[]*common.Service) *common.Service {

	var service *common.Service
	for _, s := range *services {
		if s.Name == "web" || s.Name == "custom_web" {
			service = s
			break
		}
	}
	if service == nil {
		service = &common.Service{Name: "web"}
		*services = append(*services, service)
	}
	return service
}

func (a *AnalyzerBase) AskForCommand(defaultCommand string, step string) string {
	confirmed := defaultCommand != "" && common.AskYesOrNo(fmt.Sprintf("This command will be run after each %s: '%s', confirm?", step, defaultCommand), true, a.ShouldPrompt)
	if confirmed {
		return defaultCommand
	} else {
		return common.AskUserWithDefault(fmt.Sprintf("Enter command to run after each %s:", step), "", a.ShouldPrompt)
	}
}

func (a *AnalyzerBase) refineServices(services *[]*common.Service) {
	var err error
	for _, service := range *services {
		if service.Command, err = common.ParseEnvironmentVariables(service.Command); err != nil {
			common.PrintlnError("Failed to replace environment variable placeholder due to %s", err.Error())
		}

		if service.Command, err = common.ParseUniqueInt(service.Command); err != nil {
			common.PrintlnError("Failed to replace UNIQUE_INT variable placeholder due to %s", err.Error())
		}
	}
}

func (a *AnalyzerBase) inheritProjectContext(services *[]*common.Service, envVars []*common.EnvVar, gitBranch string, gitURL string, buildRoot string) {
	for _, service := range *services {
		if service.EnvVars == nil {
			service.EnvVars = envVars
		}
		service.GitBranch = gitBranch
		service.GitRepo = gitURL
		service.BuildRoot = buildRoot
	}
}
