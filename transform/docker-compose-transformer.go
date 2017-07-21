package transform

import (
	"github.com/cloud66/starter/definitions/docker-compose"
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/service-yml"
)

type DockerComposeTransformer struct {
	Base docker_compose.DockerCompose
}

func (s *DockerComposeTransformer) ToKubernetes() kubernetes.Kubernetes {
	return kubernetes.Kubernetes{}
}

func (s *DockerComposeTransformer) ToServiceYml() service_yml.ServiceYml {
	return service_yml.ServiceYml{}
}

func (s *DockerComposeTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return docker_compose.DockerCompose{}
}
