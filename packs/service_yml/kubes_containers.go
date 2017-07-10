package service_yml

type Containers struct{
	Name string `yaml:"name,omitempty"`
	Ports []KubesPorts `yaml:"ports,omitempty"`
	Command []string `yaml:"command,omitempty"`
	Image string `yaml:"image,omitempty"`
	Resources KubesResources `yaml:"resources,omitempty"`
	Env []EnvVar `yaml:"env,omitempty"`
	VolumeMounts []VolumeMounts `"volumeMounts,omitempty"`
}

