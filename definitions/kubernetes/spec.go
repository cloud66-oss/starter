package kubernetes

type Spec struct {
	Ports []Port `yaml:"ports,omitempty"`
	Template Template `yaml:"template,omitempty"`
	Type string `yaml:"type,omitempty"`
}

type PodSpec struct{
	Containers []Containers `yaml:"containers,omitempty"`
	TerminationGracePeriodSeconds int `yaml:"terminationGracePeriodSeconds,omitempty"`
}
