package service_yml

type Limits struct{
	Cpu int `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}
