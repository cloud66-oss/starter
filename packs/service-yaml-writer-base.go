package packs

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cloud66/starter/common"
)

type ServiceYAMLWriterBase struct {
	PackElement
	TemplateDir     string
	OutputDir       string
	ShouldOverwrite bool
}

func (w *ServiceYAMLWriterBase) Write(context interface{}) error {
	templateName := fmt.Sprintf("%s.service.yml.template", w.GetPack().Name())
	if !common.FileExists(filepath.Join(w.TemplateDir, templateName)) {
		templateName = "service.yml.template" // fall back on generic template
	}

	destFullPath := filepath.Join(w.OutputDir, "service.yml")

	tmpl, err := template.ParseFiles(filepath.Join(w.TemplateDir, templateName))
	if err != nil {
		return err
	}

	if _, err := os.Stat(destFullPath); !os.IsNotExist(err) && !w.ShouldOverwrite {
		return fmt.Errorf("service.yml exists and will not be overwritten unless the overwrite flag (-o) is set")
	}

	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("%s Cannot close file service.yml due to %s\n", common.MsgError, err.Error())
		}
	}()

	fmt.Println(common.MsgL1, "Writing service.yml...")
	err = tmpl.Execute(destFile, context)
	if err != nil {
		return err
	}

	return nil
}
