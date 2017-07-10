package service_yml

type ServiceYml struct{
	Services map[string]ServiceYMLService `yaml:"services,omitempty"`
	Dbs []string `yaml:"dbs,omitempty"`
}
