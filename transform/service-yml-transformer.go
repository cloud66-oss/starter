package transform

import (
	"github.com/cloud66/starter/definitions/service-yml"
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/docker-compose"
)

type ServiceYmlTransformer struct {
	Base service_yml.ServiceYml
}

func (s *ServiceYmlTransformer) ToKubernetes() kubernetes.Kubernetes {
	return kubernetes.Kubernetes{}
}

func (s *ServiceYmlTransformer) ToServiceYml() service_yml.ServiceYml {
	return service_yml.ServiceYml{}
}

func (s *ServiceYmlTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return docker_compose.DockerCompose{}
}
