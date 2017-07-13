package service_yml

type ServiceYMLService struct {
	Image            string `yaml:"image,omitempty"`
	Volumes          []string `yaml:"volumes,omitempty"`
	StopGrace        string `yaml:"stop_grace,omitempty"`
	Constraints      Constraints  `yaml:"constraints,omitempty"`
	WorkDir          string `yaml:"work_dir,omitempty"`
	Privileged       bool `yaml:"privileged,omitempty"`
	Tags             map[string]string `yaml:"tags,omitempty"`
	Command          Command `yaml:"command,omitempty"`
	EnvVars          map[string]string `yaml:"env_vars,omitempty"`
	Ports            []interface{} `yaml:"ports,omitempty"`
	PreStopCommand   string `yaml:"pre_stop_command"`
	PostStartCommand string `yaml:"post_start_command"`

	//add unsupported keys
	GitUrl          GitUrl `yaml:"git_url,omitempty"`
	GitBranch       GitBranch `yaml:"git_branch,omitempty"`
	DockerfilePath  DockerfilePath `yaml:"dockerfile_path,omitempty"`
	Requires        Requires `yaml:"requires,omitempty"`
	BuildCommand    BuildCommand `yaml:"build_command,omitempty"`
	BuildRoot       BuildRoot `yaml:"build_root,omitempty"`
	LogFolder       LogFolder `yaml:"log_folder,omitempty"`
	DnsBehaviour    DnsBehaviour `yaml:"dns_behaviour,omitempty"`
	UseHabitus      UseHabitus `yaml:"use_habitus,omitempty"`
	UseHabitusStep  UseHabitusStep `yaml:"use_habitus_step,omitempty"`
	Health          Health `yaml:"health,omitempty"`
	PreStartSignal  PreStartSignal `yaml:"pre_start_signal,omitempty"`
	PreStopSequence PreStopSequence `yaml:"pre_stop_sequence,omitempty"`
	RestartOnDeploy RestartOnDeploy `yaml:"restart_on_deploy,omitempty"`
	TrafficMatches  TrafficMatches `yaml:"traffic_matches,omitempty"`
}
