package transform

import (
	"github.com/cloud66-oss/starter/definitions/kubernetes"
	"github.com/cloud66-oss/starter/definitions/service-yml"
)

type Transformer interface {
	ToServiceYml() service_yml.ServiceYml
	ToKubernetes() kubernetes.Kubernetes
}
