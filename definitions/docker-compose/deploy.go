package docker_compose

type Deploy struct{
	Replicas int `yaml:"replicas,omitempty"`
	UpdateConfig UpdateConfig `yaml:"update_config,omitempty"`
	RestartPolicy RestartPolicy `yaml:"restart_policy,omitempty"`
	Mode string `yaml:"mode,omitempty"`
	Placement Placement `yaml:"placement,omitempty"`
	Resources Resources `yaml:"resources,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}
