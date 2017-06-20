package transformer


type docker_Service struct {
	Command           Command `yaml:"command,omitempty"`
	Ports             Ports `yaml:"ports,omitempty"`
	Build_Command     Build_Command `yaml:"build,omitempty"`
	Image             string `yaml:"image,omitempty"`
	Depends_on        []string `yaml:"depends_on,omitempty"`
	EnvVars           []string `yaml:"environment,omitempty"`
	Deploy            Deploy `yaml:"deploy,omitempty"`
	Volumes           Volumes `yaml:"volumes,omitempty"`
	Stop_grace_period string `yaml:"stop_grace_period,omitempty"`
	Working_dir       string `yaml:"working_dir,omitempty"`
	Privileged        bool `yaml:"privileged,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty"`
	Expose            []string `yaml:"expose,omitempty"`
	Env_file          Env_file `yaml:"env_file,omitempty"`

	//unsupported docker-compose specifications
	Links          Links
	Cap_add        Cap_add
	Cap_drop       Cap_drop
	Logging        Logging
	Cgroup_parent  Cgroup_parent
	Container_name Container_name
	Devices        Devices
	Dns            Dns
	Dns_search     Dns
	External_links Links
	Extra_hosts    Extra_hosts
	Isolation      Isolation
	Networks       Networks
	Pid            Extra_hosts
	Secrets        Secrets
	Security_opt   Security_opt
	Userns_mode    Userns_mode
	Ulimits        Ulimits
	Healthcheck    Healthcheck
}
