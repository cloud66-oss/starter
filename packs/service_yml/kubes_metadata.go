package service_yml

type Metadata struct {
	Name string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	Labels map[string]string `yaml:"namespace,omitempty"`
}
