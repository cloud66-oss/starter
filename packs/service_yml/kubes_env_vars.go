package service_yml

type EnvVar struct{
	Name string `yaml:"name,omitempty"`
	Value string  `yaml:"value,omitempty"`
}
