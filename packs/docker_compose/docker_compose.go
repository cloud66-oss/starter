package docker_compose


type DockerCompose struct {
	Services map[string]DockerService
	Version string
}
