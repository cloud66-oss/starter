package docker_compose


type Links struct {
	Links []string `yaml:"links,omitempty"`
}

type CapAdd struct {
	CapAdd string
}
type CapDrop struct {
	CapDrop string
}
type ContainerName struct {
	ContainerName string
}

type CgroupParent struct {
	CgroupParent string
}

type Devices struct {
	Devices string
}

type Dns struct {
	Dns string
}

type ExtraHosts struct {
	ExtraHosts []string
}

type Isolation struct {
	Isolation string
}
type Networks struct {
	Networks []string
}
type Secrets struct {
	Secrets string
}
type SecurityOpt struct {
	SecurityOpt string
}

type UsernsMode struct {
	UsernsMode string
}

type Ulimits struct {
	Ulimits string
}

type Healthcheck struct {
	Healthcheck []string
}

type Logging struct {
	Logging []string
}
