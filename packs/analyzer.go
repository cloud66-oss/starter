package packs

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type Analyzer interface {
	GetPack() Pack
	GetRootDir() string

	Init() error
	AnalyzeServices(*[]*common.Service) error
	GuessPackages() *common.Lister
	FindVersion() string
	FindDatabases() *common.Lister
	EnvVars() []*common.EnvVar
}

type AnalyzerBase struct {
	PackElement
	RootDir     string
	Environment string
	Messages    common.Lister
}

func (a *AnalyzerBase) GetRootDir() string {
	return a.RootDir
}

func AnalyzeProcfile(a Analyzer) ([]*common.Service, error) {
	services := []*common.Service{}
	procfilePath := filepath.Join(a.GetRootDir(), "Procfile")
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

func RefineServices(services *[]*common.Service, envVars []*common.EnvVar, gitBranch string, gitURL string) {
	var err error
	for _, service := range *services {
		if service.Command, err = common.ParseEnvironmentVariables(service.Command); err != nil {
			fmt.Printf("%s Failed to replace environment variable placeholder due to %s\n", common.MsgError, err.Error())
		}

		if service.Command, err = common.ParseUniqueInt(service.Command); err != nil {
			fmt.Printf("%s Failed to replace UNIQUE_INT variable placeholder due to %s\n", common.MsgError, err.Error())
		}
		service.EnvVars = envVars
	}

	for _, service := range *services {
		service.GitBranch = gitBranch
		service.GitRepo = gitURL
	}
}
