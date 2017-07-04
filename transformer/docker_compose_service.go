package transformer

type DockerService struct {
	Command           Command `yaml:"command,omitempty"`
	Ports             Ports `yaml:"ports,omitempty"`
	BuildCommand      BuildCommand `yaml:"build,omitempty"`
	Image             string `yaml:"image,omitempty"`
	Depends_on        []string `yaml:"depends_on,omitempty"`
	EnvVars           EnvVars `yaml:"environment,omitempty"`
	Deploy            Deploy `yaml:"deploy,omitempty"`
	Volumes           Volumes `yaml:"volumes,omitempty"`
	Stop_grace_period string `yaml:"stop_grace_period,omitempty"`
	Working_dir       string `yaml:"working_dir,omitempty"`
	Privileged        bool `yaml:"privileged,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty"`
	Expose            []string `yaml:"expose,omitempty"`
	EnvFile           EnvFile `yaml:"env_file,omitempty"`
	CpuShares         int `yaml:"cpu_shares,omitempty"`
	MemLimit          int `yaml:"mem_limit,omitempty"`

	//unsupported docker-compose specifications
	Links         Links
	CapAdd        CapAdd
	CapDrop       CapDrop
	Logging       Logging
	CgroupParent  CgroupParent
	ContainerName ContainerName
	Devices       Devices
	Dns           Dns
	DnsSearch     Dns
	ExternalLinks Links
	ExtraHosts    ExtraHosts
	Isolation     Isolation
	Networks      Networks
	Pid           ExtraHosts
	Secrets       Secrets
	SecurityOpt   SecurityOpt
	UsernsMode    UsernsMode
	Ulimits       Ulimits
	Healthcheck   Healthcheck
}
