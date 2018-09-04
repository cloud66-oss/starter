package service_yml

type Ports []Port

type Port struct {
	Container int `yaml:"container,omitempty"`
	Http      int `yaml:"http,omitempty"`
	Https     int `yaml:"https,omitempty"`
	Tcp       int `yaml:"tcp,omitempty"`
	Udp       int `yaml:"udp,omitempty"`
}

type tempPort struct {
	Container string `yaml:"container,omitempty"`
	Http      string `yaml:"http,omitempty"`
	Https     string `yaml:"https,omitempty"`
	Tcp       string `yaml:"tcp,omitempty"`
	Udp       string `yaml:"udp,omitempty"`
}
