package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"gopkg.in/yaml.v2"
)

var (
	flagPath         string
	flagTemplatePath string
	flagOverwrite    bool

	packList *[]packs.Pack
	context  *common.ParseContext
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

	fmt.Println(common.MsgTitle, "Cloud 66 Starter - (c) 2015 Cloud 66", common.MsgReset)

	packList = &[]packs.Pack{&packs.Ruby{WorkDir: flagPath}}

	for _, r := range *packList {
		result, err := r.Detect()
		if err != nil {
			fmt.Printf(common.MsgError, "Failed to check for %s due to %s\n", r.Name(), err.Error())
		} else {
			if result {
				fmt.Printf("%s Found %s application\n", common.MsgL0, r.Name())
			}
		}

		if result {
			// this populates the values needed to hydrate Dockerfile.template for this pack
			context, err := r.Compile()
			if err != nil {
				fmt.Printf("Failed to compile the project due to %s", err.Error())
			}

			if err := parseAndWrite(r, fmt.Sprintf("%s.dockerfile.template", r.Name()), "Dockerfile"); err != nil {
				fmt.Printf("Failed to write Dockerfile due to %s\n", err.Error())
			}

			context, err = parseProcfile(filepath.Join(flagPath, "Procfile"), context)
			if err != nil {
				fmt.Printf("Failed to parse Procfile due to %s\n", err.Error())
			}

			if err := writeServiceFile(context, r.OutputFolder()); err != nil {
				fmt.Printf("Failed to write services.yml due to %s\n", err.Error())
			}

			break
		}
	}

	fmt.Println(common.MsgTitle, "\n Done", common.MsgReset)
}

func parseProcfile(procfilePath string, context *common.ParseContext) (*common.ParseContext, error) {
	fmt.Println(common.MsgL1, "Parsing Procfile")
	procs, err := common.ParseProcfile(procfilePath)
	if err != nil {
		return nil, err
	}

	for _, proc := range procs {
		if proc.Name == "web" || proc.Name == "custom_web" {
			// assuming the first one is created by Compile
			// TODO: this is neither safe nor right. alternatives?
			context.Services[0].Command = proc.Command
		} else {
			fmt.Printf("%s ----> Found Procfile item %s\n", common.MsgL2, proc.Name)
			context.Services = append(context.Services, &common.Service{Name: proc.Name, Command: proc.Command})
		}
	}

	// Get the git info
	gitUrl := common.RemoteGitUrl(flagPath)
	gitBranch := common.LocalGitBranch(flagPath)

	for _, service := range context.Services {
		if service.Command, err = common.ParseEnvironmentVariables(service.Command); err != nil {
			fmt.Printf("Failed to replace environment variable placeholder due to %s\n", err.Error())
		}

		if service.Command, err = common.ParseUniqueInt(service.Command); err != nil {
			fmt.Printf("Failed to replace UNIQUE_INT variable placeholder due to %s\n", err.Error())
		}

		service.GitBranch = gitBranch
		service.GitRepo = gitUrl
	}

	return context, nil
}

func writeServiceFile(context *common.ParseContext, outputFolder string) error {
	destFullPath := filepath.Join(outputFolder, "service.yml")

	d, err := yaml.Marshal(&context)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	if _, err := os.Stat(destFullPath); !os.IsNotExist(err) && !flagOverwrite {
		return fmt.Errorf("service.yml exists and will not be overwritten unless the overwrite flag is set")
	}

	fmt.Println(common.MsgL1, "Writing service.yml")
	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("Cannot close file service.yml due to %s\n", err.Error())
		}
	}()

	destFile.Write(d)

	return nil
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

	fmt.Printf("%s Writing %s...%s\n", common.MsgL1, destName, common.MsgReset)
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
