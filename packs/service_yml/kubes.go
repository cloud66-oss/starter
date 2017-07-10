package service_yml

type Kubes struct{
	Services map[string]KubesService
	Deployments map[string]KubesDeployment
}
