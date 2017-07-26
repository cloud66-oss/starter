package docker_compose

type Ulimits struct {
	Nproc  Limits
	Nofile Limits
}
