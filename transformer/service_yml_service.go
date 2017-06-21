package transformer

type ServiceYMLService struct {
	Name           string `yaml:"name,omitempty"`
	GitRepo        string `yaml:"git_url,omitempty"`
	GitBranch      string `yaml:"git_branch,omitempty"`
	BuildCommand   string `yaml:"build_command,omitempty"`
	BuildRoot      string `yaml:"build_root,omitempty"`
	Image          string `yaml:"image,omitempty"`
	Requires       []string `yaml:"requires,omitempty"`
	Volumes        []string `yaml:"volumes,omitempty"`
	StopGrace      string `yaml:"stop_grace,omitempty"`
	Constraints    Constraints  `yaml:"constraints,omitempty"`
	WorkDir        string `yaml:"work_dir,omitempty"`
	Privileged     bool `yaml:"privileged,omitempty"`
	DockerfilePath string `yaml:"dockerfile_path,omitempty"`
	Tags           map[string]string `yaml:"tags,omitempty"`
	Command        []string `yaml:"command,omitempty"`
	EnvVars        map[string]string `yaml:"env_vars,omitempty"`
	Ports          []string `yaml:"ports,omitempty"`
}

