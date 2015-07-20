package packs

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type AnalyzerBase struct {
	PackElement

	RootDir     string
	Environment string

	Messages common.Lister
}

func (a *AnalyzerBase) ProjectMetadata() (string, string, string, error) {
	gitURL := common.LocalGitBranch()
	gitBranch := common.RemoteGitUrl()
	buildRoot, err := common.PathRelativeToGitRoot(a.RootDir)
	if err != nil {
		return "", "", "", err
	}

	return gitURL, gitBranch, buildRoot, nil
}

func (b *AnalyzerBase) AnalyzeServices(a Analyzer, envVars []*common.EnvVar, gitBranch string, gitURL string, buildRoot string) ([]*common.Service, error) {
	services, err := b.analyzeProcfile()
	if err != nil {
		fmt.Printf("%s Failed to parse Procfile due to %s\n", common.MsgError, err.Error())
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

func (a *AnalyzerBase) analyzeProcfile() ([]*common.Service, error) {
	services := []*common.Service{}
	procfilePath := filepath.Join(a.RootDir, "Procfile")
	if !common.FileExists(procfilePath) {
		return services, nil
	}

	fmt.Println(common.MsgL1, "Parsing Procfile")
	procs, err := common.ParseProcfile(procfilePath)
	if err != nil {
		return nil, err
	}

	for _, proc := range procs {
		fmt.Printf("%s ----> Found Procfile item %s\n", common.MsgL2, proc.Name)
		services = append(services, &common.Service{Name: proc.Name, Command: proc.Command})
	}
	return services, nil
}

func (a *AnalyzerBase) refineServices(services *[]*common.Service) {
	var err error
	for _, service := range *services {
		if service.Command, err = common.ParseEnvironmentVariables(service.Command); err != nil {
			fmt.Printf("%s Failed to replace environment variable placeholder due to %s\n", common.MsgError, err.Error())
		}

		if service.Command, err = common.ParseUniqueInt(service.Command); err != nil {
			fmt.Printf("%s Failed to replace UNIQUE_INT variable placeholder due to %s\n", common.MsgError, err.Error())
		}
	}
}

func (a *AnalyzerBase) inheritProjectContext(services *[]*common.Service, envVars []*common.EnvVar, gitBranch string, gitURL string, buildRoot string) {
	for _, service := range *services {
		service.EnvVars = envVars
		service.GitBranch = gitBranch
		service.GitRepo = gitURL
		service.BuildRoot = buildRoot
	}
}
