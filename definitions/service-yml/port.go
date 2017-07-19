package service_yml

type Ports []Port

type Port struct {
	Container int `yaml:"container,omitempty"`
	Http      int `yaml:"http,omitempty"`
	Https     int `yaml:"https,omitempty"`
	Tcp       int `yaml:"tcp,omitempty"`
	Udp       int `yaml:"udp,omitempty"`
}
