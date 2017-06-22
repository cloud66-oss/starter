package transformer

type ServiceYml struct{
	Services map[string]ServiceYMLService `yaml:"services,omitempty"`
	Dbs []string `yaml:"dbs,omitempty"`
}
