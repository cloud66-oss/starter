package packs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloud66/starter/common"
)

type ServiceYAMLWriterBase struct {
	PackElement
	TemplateWriterBase
}

func (w *ServiceYAMLWriterBase) Write(context interface{}) error {
	templateName := fmt.Sprintf("%s.service.yml.template", w.GetPack().Name())
	if !common.FileExists(filepath.Join(w.TemplateDir, templateName)) {
		templateName = "service.yml.template" // fall back on generic template
	}
	err := w.WriteTemplate(templateName, "service.yml", context)
	if err != nil {
		return err
	}

	return w.removeBlankLines("service.yml")
}

// NOTE
// Templates can have a lot of blank lines when some parts of it are evaluated
// empty. There is currently no way to avoid this with Go templates (See
// https://github.com/golang/go/issues/9969 for more information)
// What we do to avoid this is:
// 		* Removing all blank lines from the evaluated template
// 		* If a line contains '##NEWLINE##', replace it by a newline
func (w *TemplateWriterBase) removeBlankLines(name string) error {
	fullPath := filepath.Join(w.OutputDir, name)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	blankLinePattern := regexp.MustCompile("^[[:blank:]]*$")
	newlinePattern := regexp.MustCompile("^##NEWLINE##$")

	lines := strings.Split(string(content), "\n")
	var newContent []string
	for _, line := range lines {
		if !blankLinePattern.MatchString(line) {
			if newlinePattern.MatchString(line) {
				newContent = append(newContent, "")
			} else {
				newContent = append(newContent, line)
			}
		}
	}
	newContent = append(newContent, "") // final newline

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(newContent, "\n"))
	return err
}
