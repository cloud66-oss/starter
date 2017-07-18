package service_yml

type Containers struct {
	Name            string `yaml:"name,omitempty"`
	Command         []string `yaml:"command,omitempty"`
	Image           string `yaml:"image,omitempty"`
	Resources       KubesResources `yaml:"resources,omitempty"`
	WorkingDir      string `yaml:"workingDir,omitempty"`
	SecurityContext SecurityContext `yaml:"securityContext,omitempty"`
	VolumeMounts    []VolumeMounts `yaml:"volumeMounts,omitempty"`
	Env             []EnvVar `yaml:"env,omitempty"`
	Ports           []Ports `yaml:"ports,omitempty"`
	Lifecycle       Lifecycle `yaml:"lifecycle,omitempty"`
}
