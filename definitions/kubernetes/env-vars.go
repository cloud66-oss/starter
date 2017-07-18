package service_yml

type EnvVar struct{
	Name interface{} `yaml:"name,omitempty"`
	Value interface{}  `yaml:"value,omitempty"`
}
