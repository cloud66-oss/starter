package common

type Service struct {
	Name          string
	GitRepo       string
	GitBranch     string
	Command       string
	BuildCommand  string
	DeployCommand string
	Tags          map[string]string
	Ports         []*PortMapping
	EnvVars       []*EnvVar
	BuildRoot     string
	Databases     []Database
}
