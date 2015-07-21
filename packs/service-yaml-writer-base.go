package packs

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
)

type ServiceYAMLWriterBase struct {
	PackElement
	TemplateWriterBase
	ShouldOverwrite bool
}

func (w *ServiceYAMLWriterBase) Write(context interface{}) error {
	templateName := fmt.Sprintf("%s.service.yml.template", w.GetPack().Name())
	if !common.FileExists(filepath.Join(w.TemplateDir, templateName)) {
		templateName = "service.yml.template" // fall back on generic template
	}
	return w.WriteTemplate(templateName, "service.yml", context)
}
