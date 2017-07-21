package transform

import (
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/service-yml"
	"github.com/cloud66/starter/definitions/docker-compose"
)

type KubesTransformer struct {
	Base kubernetes.Kubernetes
}

func (s *KubesTransformer) ToKubernetes() kubernetes.Kubernetes {
	return kubernetes.Kubernetes{}
}

func (s *KubesTransformer) ToServiceYml() service_yml.ServiceYml {
	return service_yml.ServiceYml{}
}

func (s *KubesTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return docker_compose.DockerCompose{}
}

