package packs

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66-oss/starter/common"
)

type DockerComposeYAMLWriterBase struct {
	PackElement
	TemplateWriterBase
}

func (w *DockerComposeYAMLWriterBase) Write(context interface{}) error {
	templateName := fmt.Sprintf("%s.docker-compose.yml.template", w.GetPack().Name())
	if !common.FileExists(filepath.Join(w.TemplateDir, templateName)) {
		templateName = "docker-compose.yml.template" // fall back on generic template
	}
	err := w.WriteTemplate(templateName, "docker-compose.yml", context)
	if err != nil {
		return err
	}

	return w.removeBlankLines("docker-compose.yml")
}
