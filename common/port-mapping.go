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

func (p PortMapping) GetEnvironmentVariablesArray(serviceName string) map[string]string {
	var envas = make(map[string]string)
	envPrefix := strings.ToUpper(serviceName)
	envSuffix := "PORT"
	if p.Container != "" {
		envas[envPrefix+"_CONTAINER_"+envSuffix] = p.Container
	}
	if p.HTTP != "" {
		envas[envPrefix+"_HTTP_"+envSuffix] = p.HTTP
	}
	if p.HTTPS != "" {
		envas[envPrefix+"_HTTPS_"+envSuffix] = p.HTTPS
	}
	if p.TCP != "" {
		envas[envPrefix+"_TCP_"+envSuffix] = p.TCP
	}
	if p.UDP != "" {
		envas[envPrefix+"_UDP_"+envSuffix] = p.UDP
	}

	return envas
}
