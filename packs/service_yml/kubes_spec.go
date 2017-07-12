package service_yml

type Spec struct {
	Ports []KubesPorts `yaml:"ports,omitempty"`
	Template Template `yaml:"template,omitempty"`
	Type string `yaml:"type,omitempty"`
}

type PodSpec struct{
	Containers []Containers `yaml:"containers,omitempty"`
	TerminationGracePeriodSeconds int `yaml:"terminationGracePeriodSeconds,omitempty"`
}
