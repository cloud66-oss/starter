package common

type Service struct {
	Name          string `yaml:"name,omitempty"`
	GitRepo       string `yaml:"git_url,omitempty"`
	GitBranch     string `yaml:"git_branch,omitempty"`
	Command       string `yaml:"command1,omitempty"`
	BuildCommand  string `yaml:"build_command,omitempty"`
	DeployCommand string `yaml:"deploy,omitempty"`
	Ports         []*PortMapping `yaml:"ports1,omitempty"`
	EnvVars       []*EnvVar `yaml:"env_vars1,omitempty"`
	BuildRoot     string `yaml:"build_root,omitempty"`
	Databases     []Database `yaml:"dbs,omitempty"`

	//stuff to be added
	Image           string `yaml:"image,omitempty"`
	Requires        []string `yaml:"requires,omitempty"`
	Volumes         []string `yaml:"volumes,omitempty"`
	Stop_grace      string `yaml:"stop_grace,omitempty"`
	Constraints     Constraints  `yaml:"constraints,omitempty"`
	Work_dir        string `yaml:"work_dir,omitempty"`
	Privileged      bool `yaml:"privileged,omitempty"`
	Dockerfile_path string `yaml:"dockerfile_path,omitempty"`
	Tags            []string `yaml:"tags,omitempty"`
	CommandSlice    []string `yaml:"command,omitempty"`
	EnvVarsSlice    []string `yaml:"env_vars,omitempty"`
	PortsShort      []string `yaml:"ports,omitempty"`
}

type Constraints struct {
	Resources Resources `yaml:"resources,omitempty"`
}

type Resources struct{
	Memory string `yaml:"memory,omitempty"`
	Cpu string  `yaml:"cpu,omitempty"`
}
