package transform

import (
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/service-yml"
	"github.com/cloud66/starter/definitions/docker-compose"
)

type KubesTransformer struct {
	Base kubernetes.Kubernetes
}

func (k *KubesTransformer) ToKubernetes() kubernetes.Kubernetes {
	return k.Base
}

func (k *KubesTransformer) ToServiceYml() service_yml.ServiceYml {
	return service_yml.ServiceYml{}
}

func (k *KubesTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return docker_compose.DockerCompose{}
}

