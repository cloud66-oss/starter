package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/kardianos/osext"
)

var (
	flagPath         string
	flagTemplatePath string
	flagOverwrite    bool
	flagEnvironment  string
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.StringVar(&flagTemplatePath, "templates", "", "where template files are located")
	flag.BoolVar(&flagOverwrite, "o", false, "overwrite existing files")
	flag.StringVar(&flagEnvironment, "e", "production", "set project environment")
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	fmt.Println(common.MsgTitle, "Cloud 66 Starter - (c) 2015 Cloud 66", common.MsgReset)

	if flagPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("%s Unable to detect current directory path due to %s", common.MsgError, err.Error())
		}
		flagPath = pwd
	}

	if flagTemplatePath == "" {
		execDir, err := osext.Executable()
		if err != nil {
			fmt.Printf("%s Unable to detect template folder due to %s", common.MsgError, err.Error())
		}

		flagTemplatePath = filepath.Join(filepath.Dir(execDir), "templates")
	}

	fmt.Printf("%s Detecting framework for the project at %s%s\n", common.MsgTitle, flagPath, common.MsgReset)

	pack, err := Detect(flagPath)
	if err != nil {
		fmt.Printf("%s Failed to detect framework due to: %s\n", common.MsgError, err.Error())
		return
	}

	err = pack.Analyze(flagPath, flagEnvironment)
	if err != nil {
		fmt.Printf("%s Failed to analyze the project due to: %s\n", common.MsgError, err.Error())
		return
	}

	err = pack.WriteDockerfile(flagTemplatePath, flagPath, flagOverwrite)
	if err != nil {
		fmt.Printf("%s Failed to write Dockerfile due to: %s\n", common.MsgError, err.Error())
	}

	err = pack.WriteServiceYAML(flagTemplatePath, flagPath, flagOverwrite)
	if err != nil {
		fmt.Printf("%s Failed to write services.yml due to: %s\n", common.MsgError, err.Error())
	}

	if len(pack.GetMessages()) > 0 {
		fmt.Printf("%s Warnings: \n", common.MsgWarn)
		for _, m := range pack.GetMessages() {
			fmt.Printf(" %s %s\n", common.MsgWarn, m)
		}
	}

	fmt.Println(common.MsgTitle, "\n Done", common.MsgReset)
}
