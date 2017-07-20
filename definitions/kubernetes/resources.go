package kubernetes

type KubesResources struct{
	Limits Limits `yaml:"limits,omitempty"`
}
