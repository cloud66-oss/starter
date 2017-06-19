package main

type Docker_compose_service struct {
	Command    []string
	Ports      []string
	Build      string
	Image      string
	Depends_on []string
	EnvVars    []string
	Deploy     string
	Volumes    []string
}
