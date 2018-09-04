package docker_compose

type Build struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
	CacheFrom  []string          `yaml:"cache_from,omitempty"`
	Labels     []string          `yaml:"labels,omitempty"`
}

type BuildAux struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
	CacheFrom  []string          `yaml:"cache_from,omitempty"`
	Labels     []string          `yaml:"labels,omitempty"`
}
