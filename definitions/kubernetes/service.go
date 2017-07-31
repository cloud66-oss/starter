package kubernetes

type KubesService struct{
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Kind string `yaml:"kind,omitempty"`
	Metadata Metadata `yaml:"metadata,omitempty"`
	Spec Spec `yaml:"spec,omitempty"`
}
