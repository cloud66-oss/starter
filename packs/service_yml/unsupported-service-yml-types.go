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

type LogFolder struct{
	LogFolder string
}


//
type DnsBehaviour struct{
	DnsBehaviour string
}

type UseHabitus struct{
	UseHabitus bool
}

type UseHabitusStep struct{
	UseHabitusStep string
}

type Health struct{
	Health string
}

type PreStartSignal struct{
	PreStartSignal string
}

type PreStopSequence struct{
	PreStopSequence string
}

type RestartOnDeploy struct{
	RestartOnDeploy bool
}

type TrafficMatches struct{
	TrafficMatches string
}
