package service_yml

type GitUrl struct{
	GitUrl string
}

type GitBranch struct{
	GitBranch string
}

type Requires struct{
	Requires []string
}

type DockerfilePath struct{
	DockerfilePath string
}

type BuildCommand struct{
	BuildCommand string
}

type BuildRoot struct{
	BuildRoot string
}