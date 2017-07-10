package service_yml

type Ports struct{
	ShortSyntax []string
	Port []ServicePort
}

type ServicePort struct{
	Container string `yaml:"container,omitempty"`
	Tcp       string `yaml:"tcp,omitempty"`
	Http      string `yaml:"http,omitempty"`
	Https     string `yaml:"https,omitempty"`
	Udp       string `yaml:"udp,omitempty"`
}
