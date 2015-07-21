package packs

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cloud66/starter/common"
)

type TemplateWriterBase struct {
	TemplateDir     string
	OutputDir       string
	ShouldNotPrompt bool
}

func (w *TemplateWriterBase) WriteTemplate(templateName string, filename string, context interface{}) error {
	tmpl, err := template.ParseFiles(filepath.Join(w.TemplateDir, templateName))
	if err != nil {
		return err
	}

	destFullPath := filepath.Join(w.OutputDir, filename)
	if w.shouldRenameExistingFile(destFullPath) {
		newName := filename + ".old"
		err = os.Rename(destFullPath, filepath.Join(w.OutputDir, newName))
		if err != nil {
			return err
		}
		fmt.Println(common.MsgL1, fmt.Sprintf("Renaming %s to %s...", filename, newName), common.MsgReset)
	}

	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("%s Cannot close file %s due to: %s\n", common.MsgError, filename, err.Error())
		}
	}()

	fmt.Println(common.MsgL1, fmt.Sprintf("Writing %s...", filename), common.MsgReset)
	err = tmpl.Execute(destFile, context)
	if err != nil {
		return err
	}

	return nil
}

func (w *TemplateWriterBase) shouldRenameExistingFile(filename string) bool {
	if !common.FileExists(filename) {
		return false
	}
	if w.ShouldNotPrompt {
		return true
	}

	message := fmt.Sprintf(" %s cannot be written as it already exists. What to do? [o: overwrite, r: rename] ", filepath.Base(filename))
	answer := "none"
	for answer != "o" && answer != "r" {
		fmt.Print(common.MsgL1, message, common.MsgReset)
		fmt.Scanln(&answer)
	}
	return answer == "r"
}
