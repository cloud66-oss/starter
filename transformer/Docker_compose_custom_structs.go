package transformer


type BuildCommand struct {
	Build         Build `yaml:"dockerfile,omitempty"`
	BuildCommand string `yaml:"build,omitempty"`
}

type Build struct {
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type Deploy struct {
	Labels map[string]string `yaml:"labels,omitempty"`
}

type Command struct {
	Command []string `yaml:"command,omitempty"`
}

type Volumes struct {
	Volumes []string `yaml:"volumes,omitempty"`
}

type EnvFile struct {
	EnvFile []string `yaml:"env_file,omitempty"`
}

type EnvVars struct {
	EnvVars []string `yaml:"environment,omitempty"`
}

type Ports struct {
	Port        []Port `yaml:"ports,omitempty"`
	ShortSyntax []string `yaml:"shortsyntax,omitempty"`
}

type Port struct {
	Target    string `yaml:"target,omitempty"`
	Published string `yaml:"published,omitempty"`
	Protocol  string `yaml:"protocol,omitempty"`
	Mode      string `yaml:"mode,omitempty"`
}