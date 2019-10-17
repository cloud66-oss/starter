package transform

import (
	"github.com/cloud66-oss/starter/definitions/kubernetes"
	"github.com/cloud66-oss/starter/definitions/service-yml"
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
