package templates

type TemplateInterface interface {
	GetName() string
	GetDependencies() []string
}
