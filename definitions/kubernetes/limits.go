package kubernetes

type Limits struct {
	Cpu    int    `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}
