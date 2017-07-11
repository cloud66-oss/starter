package service_yml

type Spec struct {
	Ports []KubesPorts `yaml:"ports,omitempety"`
	Template Template `yaml:"template,omitempty"`
}

type PodSpec struct{
	Containers []Containers `yaml:"containers,omitempty"`
	TerminationGracePeriodSeconds int `yaml:"terminationGracePeriodSeconds,omitempty"`
}
