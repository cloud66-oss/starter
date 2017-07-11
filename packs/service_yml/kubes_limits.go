package service_yml

type Limits struct{
	Cpu string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}
