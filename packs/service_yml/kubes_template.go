package service_yml

type Template struct{
	Metadata Metadata `yaml:"metadata,omitempty"`
	Spec PodSpec `yaml:"spec,omitempty"`
}
