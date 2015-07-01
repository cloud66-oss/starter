package common

type ParseContext struct {
	Services []*Service
	Dbs      []string
	EnvVars  []*EnvVar // global environment variables
	Messages []string
}
