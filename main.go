package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/kardianos/osext"
)

var (
	flagPath        string
	flagNoPrompt    bool
	flagEnvironment string
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.BoolVar(&flagNoPrompt, "y", false, "do not prompt user")
	flag.StringVar(&flagEnvironment, "e", "production", "set project environment")
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	common.PrintlnTitle("Cloud 66 Starter ~ (c) 2015 Cloud 66")

	if flagPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			common.PrintlnError("Unable to detect current directory path due to %s", err.Error())
		}
		flagPath = pwd
	}

	execDir, err := osext.Executable()
	if err != nil {
		common.PrintlnError("Unable to detect template folder due to %s", err.Error())
	}
	dockerfileTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", "dockerfiles")
	serviceYAMLTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", "service-yml")

	common.PrintlnTitle("Detecting framework for the project at %s", flagPath)

	pack, err := Detect(flagPath)
	if err != nil {
		common.PrintlnError("Failed to detect framework due to: %s", err.Error())
		return
	}

	err = pack.Analyze(flagPath, flagEnvironment, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to analyze the project due to: %s", err.Error())
		return
	}

	err = pack.WriteDockerfile(dockerfileTemplateDir, flagPath, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to write Dockerfile due to: %s", err.Error())
	}

	err = pack.WriteServiceYAML(serviceYAMLTemplateDir, flagPath, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to write service.yml due to: %s", err.Error())
	}

	if len(pack.GetMessages()) > 0 {
		common.PrintlnWarning("Warnings:")
		for _, warning := range pack.GetMessages() {
			common.PrintlnWarning(warning)
		}
	}

	common.PrintlnTitle("Done")
}
