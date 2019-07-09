package common

import "strings"

type PortMapping struct {
	Container string
	HTTP      string
	HTTPS     string
	TCP       string
	UDP       string
}

func NewPortMapping() *PortMapping {
	return &PortMapping{HTTP: "80", HTTPS: "443"}
}

func NewInternalPortMapping(container string) *PortMapping {
	return &PortMapping{Container: container}
}
