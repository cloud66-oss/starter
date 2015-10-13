package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type downloadFile struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type templateDefinition struct {
	Version     string         `json:"version"`
	Dockerfiles []downloadFile `json:"dockerfiles"`
	ServiceYmls []downloadFile `json:"service-ymls"`
}

var (
	flagPath        string
	flagNoPrompt    bool
	flagEnvironment string
	flagTemplates   string
	VERSION         string = "dev"
	BUILD_DATE      string = ""

	serviceYAMLTemplateDir string
	dockerfileTemplateDir  string
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.BoolVar(&flagNoPrompt, "y", false, "do not prompt user")
	flag.StringVar(&flagEnvironment, "e", "production", "set project environment")
	flag.StringVar(&flagTemplates, "templates", "", "location of the templates directory")
}

// downloading templates from github and putting them into homedir
func getTempaltes(tempDir string) error {
	common.PrintL0("Checking templates in %s", tempDir)

	var tv templateDefinition
	err := fetchJSON("https://raw.githubusercontent.com/cloud66/starter/master/templates/templates.json", nil, tv)
	if err != nil {
		return err
	}

	// is there a local copy?
	if _, err := os.Stat(filepath.Join(tempDir, "templates.json")); os.IsNotExist(err) {
		// no file. downloading
		err := downloadTemplates(tempDir, tv)
		if err != nil {
			return err
		}
	}

	// load the local json
	templatesLocal, err := ioutil.ReadFile(filepath.Join(tempDir, "templates.json"))
	if err != nil {
		return err
	}
	var localTv templateDefinition
	err = json.Unmarshal(templatesLocal, &localTv)
	if err != nil {
		return err
	}

	// compare
	if localTv.Version != tv.Version {
		// they are different, we need to download the new ones
		err = downloadTemplates(tempDir, tv)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadTemplates(tempDir string, td templateDefinition) error {
	for _, temp := range td.Dockerfiles {
		r, err := fetch(temp.Url, nil)
		defer r.Close()

		if err != nil {
			return err
		}

		output, err := os.Create(temp.Name)
		if err != nil {
			return err
		}
		defer output.Close()

		_, err = io.Copy(output, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		flag.PrintDefaults()
		return
	}

	if len(args) > 0 && args[0] == "version" {
		common.PrintlnTitle("Starter version: %s (%s)", VERSION, BUILD_DATE)
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

	// if templateFolder is specified we're going to use that otherwise download
	if flagTemplates == "" {
		usr, _ := user.Current()
		homeDir := usr.HomeDir

		flagTemplates = filepath.Join(homeDir, ".starter")
		err := getTempaltes(flagTemplates)
		if err != nil {
			common.PrintlnError("Failed to download latest templates due to %s", err.Error())
			os.Exit(1)
		}

		dockerfileTemplateDir = flagTemplates
		serviceYAMLTemplateDir = flagTemplates
	} else {
		common.PrintlnTitle("Using local templates at %s", flagTemplates)
		flagTemplates, err := filepath.Abs(flagTemplates)
		if err != nil {
			common.PrintlnError("Failed to use %s for templates due to %s", flagTemplates, err.Error())
			os.Exit(1)
		}
		dockerfileTemplateDir = flagTemplates
		serviceYAMLTemplateDir = flagTemplates
	}

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
			common.PrintlnWarning(" * " + warning)
		}
	}

	common.PrintlnTitle("Done")
}
