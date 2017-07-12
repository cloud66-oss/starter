package service_yml

type KubesPorts struct {
	Name string `yaml:"name,omitempty"`
	Port string `yaml:"port,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
	TargetPort string `yaml:"targetPort,omitempty"`
	ContainerPort string `yaml:"containerPort,omitempty"`
}
