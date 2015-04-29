package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	//	"github.com/mgutz/ansi"
	"github.com/cloud66/starter/packs"
)

var (
	flagPath         string
	flagTemplatePath string
	flagOverwrite    bool

	packList *[]packs.Pack
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.StringVar(&flagTemplatePath, "templates", "templates", "where template files are located")
	flag.BoolVar(&flagOverwrite, "o", false, "overwrite existing files")
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	fmt.Println("Cloud 66 Starter - (c) 2015 Cloud 66")

	packList = &[]packs.Pack{&packs.Ruby{WorkDir: flagPath}}

	for _, r := range *packList {
		result, err := r.Detect()
		if err != nil {
			fmt.Printf("Failed to check for %s due to %s\n", r.Name(), err.Error())
		} else {
			if result {
				fmt.Printf("Found %s application\n", r.Name())
			}
		}

		if result {
			// this populates the values needed to hydrate Dockerfile.template for this pack
			r.Compile()

			if err := parseAndWrite(r, fmt.Sprintf("%s.dockerfile.template", r.Name()), "Dockerfile"); err != nil {
				fmt.Printf("Failed to write Dockerfile due to %s\n", err.Error())
			}

			if err := parseAndWrite(r, fmt.Sprintf("%s.service_yml.template", r.Name()), "services.yml"); err != nil {
				fmt.Printf("Failed to write services.yml due to %s\n", err.Error())
			}
		}
	}

	fmt.Println("\nDone")
}

func parseAndWrite(pack packs.Pack, templateName string, destName string) error {
	tmpl, err := template.ParseFiles(filepath.Join(flagTemplatePath, templateName))
	if err != nil {
		return err
	}

	destFullPath := filepath.Join(pack.OutputFolder(), destName)

	if _, err := os.Stat(destFullPath); !os.IsNotExist(err) && !flagOverwrite {
		return fmt.Errorf("File %s exists and will not be overwritten unless the overwrite flag is set\n", destName)
	}

	fmt.Printf("Writing %s...\n", destName)
	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("Cannot close file %s due to %s\n", destName, err.Error())
		}
	}()
	err = tmpl.Execute(destFile, pack)
	if err != nil {
		return err
	}

	return nil
}
