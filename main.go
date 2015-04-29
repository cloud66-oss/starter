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
	flagPath string

	packList *[]packs.Pack
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
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
			tmpl, err := template.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.dockerfile.template", r.Name())))
			if err != nil {
				panic(err)
			}
			err = tmpl.Execute(os.Stdout, r)
			if err != nil {
				panic(err)
			}

			tmpl, err = template.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.service_yml.template", r.Name())))
			if err != nil {
				panic(err)
			}
			err = tmpl.Execute(os.Stdout, r)
			if err != nil {
				panic(err)
			}

			// TODO: check if those files exist in the repo before overwritting
		}
	}

	fmt.Println("\nDone")
}
