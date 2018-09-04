package docker_compose

type Secrets []Secret

type Secret struct {
	Source string `yaml:"source,omitempty"`
	Target string `yaml:"target,omitempty"`
	Uid    string `yaml:"uid,omitempty"`
	Gid    string `yaml:"gid,omitempty"`
	Mode   string `yaml:"mode,omitempty"`
}
