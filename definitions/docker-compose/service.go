package docker_compose

type Service struct {
	Command           Command `yaml:"command,omitempty"`
	Ports             Ports `yaml:"ports,omitempty"`
	Build             Build `yaml:"build,omitempty"`
	Image             string `yaml:"image,omitempty"`
	Depends_on        []string `yaml:"depends_on,omitempty"`
	Environment       Environment `yaml:"environment,omitempty"`
	Deploy            Deploy `yaml:"deploy,omitempty"`
	Volumes           Volumes `yaml:"volumes,omitempty"`
	Stop_grace_period string `yaml:"stop_grace_period,omitempty"`
	Working_dir       string `yaml:"working_dir,omitempty"`
	Privileged        bool `yaml:"privileged,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty"`
	Expose            []string `yaml:"expose,omitempty"`
	EnvFile           EnvFile `yaml:"env_file,omitempty"`
	CpuShares         int `yaml:"cpu_shares,omitempty"`
	MemLimit          string `yaml:"mem_limit,omitempty"`
	Links             []string `yaml:"links,omitempty"`
	CapAdd            []string `yaml:"cap_add,omitempty"`
	CapDrop           []string `yaml:"cap_drop,omitempty"`
	Logging           Logging `yaml:"logging,omitempty"`
	CgroupParent      string `yaml:"cgroup_parent,omitempty"`
	ContainerName     string `yaml:"container_name,omitempty"`
	Devices           []string `yaml:"devices,omitempty"`
	Dns               Dns `yaml:"dns,omitempty"`
	DnsSearch         DnsSearch `yaml:"dns_search,omitempty"`
	Entrypoint        Entrypoint `yaml:"entrypoint,omitempty"`
	ExternalLinks     []string `yaml:"external_links,omitempty"`
	ExtraHosts        []string `yaml:"extra_hosts,omitempty"`
	Isolation         string `yaml:"isolation,omitempty"`
	Networks          Networks `yaml:"networks,omitempty"`
	Pid               string `yaml:"pid,omitempty"`
	Secrets           Secrets `yaml:"secrets,omitempty"`
	SecurityOpt       []string `yaml:"security_opt,omitempty"`
	StopSignal        string `yaml:"stop_signal,omitempty"`
	Sysctls           map[string]string `yaml:"sysctls,omitempty"`
	Tmpfs             Tmpfs `yaml:"tmpfs,omitempty"`
	UsernsMode        string `yaml:"userns_mode,omitempty"`
	Ulimits           Ulimits `yaml:"ulimits,omitempty"`
	Healthcheck       Healthcheck `yaml:"healthcheck,omitempty"`
}
