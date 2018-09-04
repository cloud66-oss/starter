package transform

import (
	"github.com/cloud66-oss/starter/definitions/docker-compose"
	"github.com/cloud66-oss/starter/definitions/kubernetes"
	"github.com/cloud66-oss/starter/definitions/service-yml"
)

type Transformer interface {
	ToServiceYml() service_yml.ServiceYml
	ToDockerCompose() docker_compose.DockerCompose
	ToKubernetes() kubernetes.Kubernetes
}
