package packs

import "fmt"

type DockerfileWriterBase struct {
	PackElement
	TemplateWriterBase
	ShouldOverwrite bool
}

func (w *DockerfileWriterBase) Write(context interface{}) error {
	templateName := fmt.Sprintf("%s.dockerfile.template", w.GetPack().Name())
	return w.WriteTemplate(templateName, "Dockerfile", context)
}
