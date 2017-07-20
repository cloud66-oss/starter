package kubernetes

type Ports struct {
	Name string `yaml:"name,omitempty"`
	Port int `yaml:"port,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
	TargetPort int `yaml:"targetPort,omitempty"`
	ContainerPort int `yaml:"containerPort,omitempty"`
	NodePort int `yaml:"nodePort,omitempty"`
}
