package docker_compose

type Ports []Port

type Port struct {
	Target    int    `yaml:"target,omitempty"`
	Published int    `yaml:"published,omitempty"`
	Protocol  string `yaml:"protocol,omitempty"`
	Mode      string `yaml:"mode,omitempty"`
}
