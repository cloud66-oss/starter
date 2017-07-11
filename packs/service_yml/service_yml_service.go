package service_yml

type ServiceYMLService struct {
	Name           string `yaml:"name,omitempty"`
	Image          string `yaml:"image,omitempty"`
	Volumes        []interface{} `yaml:"volumes,omitempty"`
	StopGrace      string `yaml:"stop_grace,omitempty"`
	Constraints    Constraints  `yaml:"constraints,omitempty"`
	WorkDir        string `yaml:"work_dir,omitempty"`
	Privileged     bool `yaml:"privileged,omitempty"`
	Tags           map[string]string `yaml:"tags,omitempty"`
	Command        string `yaml:"command,omitempty"`
	EnvVars        map[string]string `yaml:"env_vars,omitempty"`
	Ports          []interface{} `yaml:"ports,omitempty"`

	//add unsupported keys
	GitUrl        GitUrl `yaml:"git_url,omitempty"`
	GitBranch      GitBranch `yaml:"git_branch,omitempty"`
	DockerfilePath DockerfilePath `yaml:"dockerfile_path,omitempty"`
	Requires       Requires `yaml:"requires,omitempty"`
	BuildCommand   BuildCommand `yaml:"build_command,omitempty"`
	BuildRoot      BuildRoot `yaml:"build_root,omitempty"`
}
