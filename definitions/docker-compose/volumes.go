package docker_compose

type Volumes []Volume

type Volume struct {
	Type     string `yaml:"type,omitempty"`
	Source   string `yaml:"source,omitempty"`
	Target   string `yaml:"target,omitempty"`
	ReadOnly bool `yaml:"read_only,omitempty"`
}
