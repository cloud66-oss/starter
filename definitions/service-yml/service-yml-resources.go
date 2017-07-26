package service_yml

type Resources struct {
	Memory string `yaml:"memory,omitempty"`
	Cpu    int  `yaml:"cpu,omitempty"`
}
