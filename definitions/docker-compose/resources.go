package docker_compose

type Resources struct {
	Limits       CpusMem `yaml:"limits,omitempty"`
	Reservations CpusMem `yaml:reservations",omitempty"`
}
