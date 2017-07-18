package service_yml

type Metadata struct {
	Name string `yaml:"name,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}
