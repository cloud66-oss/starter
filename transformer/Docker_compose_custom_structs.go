package transformer


type Build_Command struct {
	Build         Build `yaml:"dockerfile,omitempty"`
	Build_Command string `yaml:"build,omitempty"`
}

type Build struct {
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type Deploy struct {
	Resources Resources `yaml:"resources,omitempty"`
}

type Resources struct {
	Limits Limits `yaml:"limits,omitempty"`
}

type Limits struct {
	Cpus   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type Command struct {
	Command []string `yaml:"command,omitempty"`
}

type Volumes struct {
	Volumes []string `yaml:"volumes,omitempty"`
}

type Env_file struct {
	Env_file []string `yaml:"env_file,omitempty"`
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