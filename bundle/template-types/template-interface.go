package template_types

type TemplateInterface interface {
	GetName() string
	GetDependencies() []string
}
