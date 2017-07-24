package kubernetes

type Resources struct{
	Limits Limits `yaml:"limits,omitempty"`
}
