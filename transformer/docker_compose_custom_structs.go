package transformer


type BuildCommand struct {
	Build         Build `yaml:"dockerfile,omitempty"`
	BuildCommand string `yaml:"build,omitempty"`
}

type Build struct {
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type Deploy struct {
	Resources DockerResources`yaml:"resources,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

type DockerResources struct{
	Limits CpusMem `yaml:"limits,omitempty"`
	Reservations CpusMem `yaml:reservations",omitempty"`

}

type CpusMem struct{
	Cpus string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type Command struct {
	Command []string `yaml:"command,omitempty"`
}

type Volumes struct {
	Volumes []string `yaml:"volumes,omitempty"`
	LongSyntax	[]LongSyntaxVolume
}

type LongSyntaxVolume struct{
	Type	string `yaml:"type,omitempty"`
	Source	string `yaml:"source,omitempty"`
	Target	string `yaml:"target,omitempty"`
	ReadOnly bool `yaml:"read_only,omitempty"`
}

type EnvFile struct {
	EnvFile []string `yaml:"env_file,omitempty"`
}

type EnvVars struct {
	EnvVars map[string]string `yaml:"environment,omitempty"`
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