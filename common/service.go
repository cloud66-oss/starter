package common

type Service struct {
	Name      string
	GitRepo   string
	GitBranch string
	Command   string
	Ports     []string
}
