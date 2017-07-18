package service_yml

type KubesResources struct{
	Limits Limits `yaml:"limits,omitempty"`
}
