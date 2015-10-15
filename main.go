package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-go"
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
	flagBranch      string
	flagOverwrite   bool
	VERSION         string = "dev"
	BUILD_DATE      string = ""

	serviceYAMLTemplateDir string
	dockerfileTemplateDir  string
)

const (
	templatePath = "https://raw.githubusercontent.com/cloud66/starter/{{.branch}}/templates/templates.json"
)

func init() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:     "916591d12b54e689edde67e641c5843d",
		AppVersion: VERSION,
	})

	flag.StringVar(&flagPath, "p", "", "project path")
	flag.BoolVar(&flagNoPrompt, "y", false, "do not prompt user")
	flag.BoolVar(&flagOverwrite, "overwrite", false, "overwrite existing files")
	flag.StringVar(&flagEnvironment, "e", "production", "set project environment")
	flag.StringVar(&flagTemplates, "templates", "", "location of the templates directory")
	flag.StringVar(&flagBranch, "branch", "master", "template branch in github")
}

// downloading templates from github and putting them into homedir
func getTempaltes(tempDir string) error {
	common.PrintlnL0("Checking templates in %s", tempDir)

	var tv templateDefinition
	err := fetchJSON(strings.Replace(templatePath, "{{.branch}}", flagBranch, -1), nil, &tv)
	if err != nil {
		return err
	}

	// is there a local copy?
	if _, err := os.Stat(filepath.Join(tempDir, "templates.json")); os.IsNotExist(err) {
		// no file. downloading
		common.PrintlnL1("No local templates found. Downloading now.")
		err := os.MkdirAll(tempDir, 0777)
		if err != nil {
			return err
		}

		err = downloadTemplates(tempDir, tv)
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
		common.PrintlnL2("Newer templates found. Downloading them now")
		// they are different, we need to download the new ones
		err = downloadTemplates(tempDir, tv)
		if err != nil {
			return err
		}
	} else {
		common.PrintlnL1("Local templates are up to date")
	}

	return nil
}

func downloadTemplates(tempDir string, td templateDefinition) error {
	err := downloadSingleFile(tempDir, downloadFile{Url: strings.Replace(templatePath, "{{.branch}}", flagBranch, -1), Name: "templates.json"})
	if err != nil {
		return err
	}

	for _, temp := range td.Dockerfiles {
		err := downloadSingleFile(tempDir, temp)
		if err != nil {
			return err
		}
	}

	for _, temp := range td.ServiceYmls {
		err := downloadSingleFile(tempDir, temp)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadSingleFile(tempDir string, temp downloadFile) error {
	r, err := fetch(strings.Replace(temp.Url, "{{.branch}}", flagBranch, -1), nil)
	if err != nil {
		return err
	}
	defer r.Close()

	output, err := os.Create(filepath.Join(tempDir, temp.Name))
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, r)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		flag.PrintDefaults()
		return
	}

	if len(args) > 0 && (args[0] == "version" || args[0] == "v") {
		fmt.Printf("Cloud 66 Starter (%s)\n", VERSION)
		fmt.Println("Copyright 2015 Cloud66 Inc.")
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
		os.Exit(2)
	}

	// check for Dockerfile (before analysis to avoid wasting time)
	dockerfilePath := filepath.Join(flagPath, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err == nil {
		// file exists. should we overwrite?
		if !flagOverwrite {
			common.PrintlnError("Dockerfile already exists. Use overwrite flag to overwrite it")
			os.Exit(1)
		}
	}

	serviceYAMLPath := filepath.Join(flagPath, "service.yml")
	if _, err := os.Stat(serviceYAMLPath); err == nil {
		// file exists. should we overwrite?
		if !flagOverwrite {
			common.PrintlnError("service.yml already exists. Use overwrite flag to overwrite it")
			os.Exit(1)
		}
	}

	err = pack.Analyze(flagPath, flagEnvironment, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to analyze the project due to: %s", err.Error())
		os.Exit(2)
	}

	err = pack.WriteDockerfile(dockerfileTemplateDir, flagPath, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to write Dockerfile due to: %s", err.Error())
		os.Exit(2)
	}

	err = pack.WriteServiceYAML(serviceYAMLTemplateDir, flagPath, !flagNoPrompt)
	if err != nil {
		common.PrintlnError("Failed to write service.yml due to: %s", err.Error())
		os.Exit(2)
	}

	if len(pack.GetMessages()) > 0 {
		common.PrintlnWarning("Warnings:")
		for _, warning := range pack.GetMessages() {
			common.PrintlnWarning(" * " + warning)
		}
	}

	common.PrintlnL0("Now you can add the newly created Dockerfile to your git and create a stack")
	common.PrintlnL0("To do that you will need to run the following commands:\n\n")
	fmt.Printf("cd %s\n", flagPath)
	fmt.Println("git add Dockerfile")
	fmt.Println("git commit -m 'Adding Dockerfile'")
	fmt.Printf("cx stacks create --name='CHANGEME' --environment='%s' --service_yaml=service.yml\n\n",
		flagEnvironment)

	common.PrintlnTitle("Done")
}
