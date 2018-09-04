package kubernetes

type Handler struct {
	Exec Exec `yaml:"exec,omitempty"`
}
