package service_yml

type KubesPorts struct {
	Port string `yaml:"port,omitempty"`
	TargetPort string `yaml:"targetPort,omitempty"`
	ContainerPort string `yaml:"containerPort,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
	NodePort string `yaml:"nodePort,omitempty"`
	Name string `yaml:"name,omitempty"`
}
