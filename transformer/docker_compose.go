package transformer


type Docker_compose struct {
	Services map[string]docker_Service
	Version string
}

type Service_yml struct {
	Services map[string]ServiceYMLService
}

type Serviceyml struct{
	Services map[string]ServiceYMLService `yaml:"services,omitempty"`
	Dbs []string `yaml:"dbs,omitempty"`
}