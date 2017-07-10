package service_yml

type Ports struct {
	Port int `yaml:"port,omitempty"`
	TargetPort int `yaml:"targetPort,omitempty"`
	ContainerPort int `yaml:"containerPort,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
	NodePort int `yaml:"nodePort,omitempty"`
	Name string `yaml:"name,omitempty"`
}
