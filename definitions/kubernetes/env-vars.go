package kubernetes

type EnvVar struct{
	Name interface{} `yaml:"name,omitempty"`
	Value interface{}  `yaml:"value,omitempty"`
}
