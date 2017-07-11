package service_yml

type Containers struct{
	Name string `yaml:"name,omitempty"`
	Ports []KubesPorts `yaml:"ports,omitempty"`
	Command string `yaml:"command,omitempty"`
	Image string `yaml:"image,omitempty"`
	Resources KubesResources `yaml:"resources,omitempty"`
	Env []EnvVar `yaml:"env,omitempty"`
	VolumeMounts []VolumeMounts `yaml:"volumeMounts,omitempty"`
	WorkingDir string `yaml:"workingDir,omitempty"`
	SecurityContext SecurityContext `yaml:"securityContext,omitempty"`
}

