package service_yml

type Lifecycle struct{
	PostStart Handler `yaml:"postStart,omitempty"`
	PreStop Handler `yaml:"preStop,omitempty"`
}
