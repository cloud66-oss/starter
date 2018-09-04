package kubernetes

type Template struct {
	Metadata Metadata `yaml:"metadata,omitempty"`
	PodSpec  PodSpec  `yaml:"spec,omitempty"`
}
