package packs

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type Analyzer interface {
	Name() string
	GetRootDir() string
	GetContext() *common.ParseContext
	GetGitBranch() string
	GetGitURL() string
	GetVersion() string
	GetMessages() common.Lister
	GetPackages() *common.Lister
	Analyze() error
	FetchGitMetadata()
	AnalyzeProcfile() error
}

type AnalyzerBase struct {
	RootDir     string
	Environment string
	Version     string

	GitUrl    string
	GitBranch string

	Packages *common.Lister
	Messages common.Lister
	Context  *common.ParseContext
}

func (a *AnalyzerBase) GetMessages() common.Lister {
	return a.Messages
}

func (a *AnalyzerBase) GetPackages() *common.Lister {
	return a.Packages
}

func (a *AnalyzerBase) GetRootDir() string {
	return a.RootDir
}

func (a *AnalyzerBase) GetVersion() string {
	return a.Version
}

func (a *AnalyzerBase) GetContext() *common.ParseContext {
	return a.Context
}

func (a *AnalyzerBase) GetGitBranch() string {
	return a.GitBranch
}

func (a *AnalyzerBase) GetGitURL() string {
	return a.GitUrl
}

func Analyze(a Analyzer) error {
	a.FetchGitMetadata()

	err := a.Analyze()
	if err != nil {
		return err
	}

	err = a.AnalyzeProcfile()
	if err != nil {
		fmt.Printf("%s Failed to parse Procfile due to %s\n", common.MsgError, err.Error())
	}

	for _, service := range a.GetContext().Services {
		if service.Command, err = common.ParseEnvironmentVariables(service.Command); err != nil {
			fmt.Printf("%s Failed to replace environment variable placeholder due to %s\n", common.MsgError, err.Error())
		}

		if service.Command, err = common.ParseUniqueInt(service.Command); err != nil {
			fmt.Printf("%s Failed to replace UNIQUE_INT variable placeholder due to %s\n", common.MsgError, err.Error())
		}

		service.EnvVars = a.GetContext().EnvVars
	}

	for _, service := range a.GetContext().Services {
		service.GitBranch = a.GetGitBranch()
		service.GitRepo = a.GetGitURL()
	}

	return nil
}

func (a *AnalyzerBase) FetchGitMetadata() {
	a.GitBranch = common.LocalGitBranch(a.RootDir)
	a.GitUrl = common.RemoteGitUrl(a.RootDir)
}

func (a *AnalyzerBase) AnalyzeProcfile() error {
	procfilePath := filepath.Join(a.RootDir, "Procfile")
	if !common.FileExists(procfilePath) {
		return nil
	}

	fmt.Println(common.MsgL1, "Parsing Procfile")
	procs, err := common.ParseProcfile(procfilePath)
	if err != nil {
		return err
	}

	for _, proc := range procs {
		if proc.Name == "web" || proc.Name == "custom_web" {
			a.Context.Services[0].Command = proc.Command
		} else {
			fmt.Printf("%s ----> Found Procfile item %s\n", common.MsgL2, proc.Name)
			a.Context.Services = append(a.Context.Services, &common.Service{Name: proc.Name, Command: proc.Command})
		}
	}
	return nil
}
