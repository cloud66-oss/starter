package service_yml

type Service struct {
	GitUrl           string            `yaml:"git_url,omitempty"`
	GitBranch        string            `yaml:"git_branch,omitempty"`
	BuildRoot        string            `yaml:"build_root,omitempty"`
	BuildCommand     string            `yaml:"build_command,omitempty"`
	DockerfilePath   string            `yaml:"dockerfile_path,omitempty"`
	Image            string            `yaml:"image,omitempty"`
	Command          string            `yaml:"command,omitempty"`
	Requires         []string          `yaml:"requires,omitempty"`
	Tags             map[string]string `yaml:"tags,omitempty"`
	Ports            Ports             `yaml:"ports,omitempty"`
	EnvVars          map[string]string `yaml:"env_vars,omitempty"`
	Volumes          []string          `yaml:"volumes,omitempty"`
	StopGrace        int               `yaml:"stop_grace,omitempty"`
	Constraints      Constraints       `yaml:"constraints,omitempty"`
	WorkDir          string            `yaml:"work_dir,omitempty"`
	Privileged       bool              `yaml:"privileged,omitempty"`
	PreStopCommand   string            `yaml:"pre_stop_command,omitempty"`
	PostStartCommand string            `yaml:"post_start_command,omitempty"`
	DeployCommand    string            `yaml:"deploy_command,omitempty"`
	LogFolder        string            `yaml:"log_folder,omitempty"`
	DnsBehaviour     string            `yaml:"dns_behaviour,omitempty"`
	UseHabitus       bool              `yaml:"use_habitus,omitempty"`
	UseHabitusStep   string            `yaml:"use_habitus_step,omitempty"`
	Health           string            `yaml:"health,omitempty"`
	PreStartSignal   string            `yaml:"pre_start_signal,omitempty"`
	PreStopSequence  string            `yaml:"pre_stop_sequence,omitempty"`
	RestartOnDeploy  bool              `yaml:"restart_on_deploy,omitempty"`
	TrafficMatches   TrafficMatches    `yaml:"traffic_matches,omitempty"`
}
