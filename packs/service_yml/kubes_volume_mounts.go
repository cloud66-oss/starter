package service_yml

type VolumeMounts struct {
	MountPath string `yaml:"mountPath,omitempty"`
	Name string `yaml:"name,omitempty"`
	ReadOnly bool `yaml:"readOnly,omitempty"`
}
