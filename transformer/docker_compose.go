package transformer


type DockerCompose struct {
	Services map[string]DockerService
	Version string
}
