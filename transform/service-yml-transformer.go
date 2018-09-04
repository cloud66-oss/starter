package transform

import (
	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/definitions/docker-compose"
	"github.com/cloud66-oss/starter/definitions/kubernetes"
	"github.com/cloud66-oss/starter/definitions/service-yml"
)

type ServiceYmlTransformer struct {
	Base service_yml.ServiceYml
}

func (s *ServiceYmlTransformer) ToKubernetes() kubernetes.Kubernetes {

	var deploy kubernetes.KubesDeployment
	var deployments []kubernetes.KubesDeployment
	var kubesServices []kubernetes.KubesService

	//Each service needs an unique nodePort, so we hand-pick to start
	//from 31111 and pray that it will not collide with other stuff.
	nodePort := 31111

	for _, dbName := range s.Base.Databases {
		if dbName != "" {
			tags := make(map[string]string, 1)
			tags["app"] = dbName

			service := kubernetes.KubesService{
				ApiVersion: "v1",
				Kind:       "Service",
				Metadata: kubernetes.Metadata{
					Name: dbName + "-svc",
				},
				Spec: kubernetes.Spec{
					Type:  "ClusterIP",
					Ports: setDbServicePorts(dbName),
				},
			}
			deploy := kubernetes.KubesDeployment{ApiVersion: "extensions/v1beta1",
				Kind: "Deployment",
				Metadata: kubernetes.Metadata{
					Name: dbName + "-deployment",
				},
				Spec: kubernetes.Spec{
					Template: kubernetes.Template{
						Metadata: kubernetes.Metadata{
							Labels: tags,
						},
						PodSpec: kubernetes.PodSpec{
							Containers: []kubernetes.Containers{
								{
									Name:  dbName,
									Image: getDbImage(dbName),
									Ports: setDbDeploymentPorts(dbName),
								},
							},
						},
					},
				},
			}

			kubesServices = append(kubesServices, service)
			deployments = append(deployments, deploy)
		}
	}
	var deployPorts []kubernetes.Port
	var services []kubernetes.KubesService
	for serviceName, serviceSpecs := range s.Base.Services {
		getServiceToKubesWarnings(serviceSpecs)
		//gets ports to populate deployment and generates the required service(s)
		deployPorts, services, nodePort = generateServicesRequiredByPorts(serviceName, serviceSpecs, nodePort)

		//required by the kubes format
		if serviceSpecs.Tags == nil {
			serviceSpecs.Tags = make(map[string]string, 1)
		}
		serviceSpecs.Tags["app"] = serviceName

		deploy = kubernetes.KubesDeployment{
			ApiVersion: "extensions/v1beta1",
			Kind:       "Deployment",
			Metadata: kubernetes.Metadata{
				Name: serviceName + "-deployment",
			},
			Spec: kubernetes.Spec{
				Template: kubernetes.Template{
					Metadata: kubernetes.Metadata{
						Labels: serviceSpecs.Tags,
					},
					PodSpec: kubernetes.PodSpec{
						TerminationGracePeriodSeconds: serviceSpecs.StopGrace,
						Containers: []kubernetes.Containers{
							{
								Name:       serviceName,
								Image:      serviceSpecs.Image,
								Command:    serviceToKubesCommand(serviceSpecs.Command),
								Ports:      deployPorts,
								WorkingDir: serviceSpecs.WorkDir,
								SecurityContext: kubernetes.SecurityContext{
									Priviliged: serviceSpecs.Privileged,
								},
								Lifecycle: kubernetes.Lifecycle{
									PostStart: kubernetes.Handler{
										Exec: kubernetes.Exec{
											Command: serviceToKubesCommand(serviceSpecs.PostStartCommand),
										},
									},
									PreStop: kubernetes.Handler{
										Exec: kubernetes.Exec{
											Command: serviceToKubesCommand(serviceSpecs.PreStopCommand),
										},
									},
								},
								Resources: kubernetes.Resources{
									Limits: kubernetes.Limits{
										Cpu:    serviceSpecs.Constraints.Resources.Cpu,
										Memory: serviceSpecs.Constraints.Resources.Memory,
									},
								},
							},
						},
					},
				},
			},
		}

		//if it has no image, output warning to user about the fact that each container needs one
		if serviceSpecs.Image == "" {
			deploy.Spec.Template.PodSpec.Containers[0].Image = "#INSERT REQUIRED IMAGE"
			common.PrintlnWarning("The service \"%s\" has no image mentioned and each container needs one in Kubernetes format. Please add manually.", serviceName)
		}

		kubeVolumes := handleVolumes(serviceSpecs.Volumes)
		deploy.Spec.Template.PodSpec.Containers[0].VolumeMounts = kubeVolumes

		keys, values := getKeysValues(serviceSpecs.EnvVars)
		if len(keys) > 0 {
			for k := 0; k < len(keys); k++ {
				if values[k] == "\"\"" {
					values[k] = ""
				}
				env := kubernetes.EnvVar{
					Name:  keys[k],
					Value: values[k],
				}
				deploy.Spec.Template.PodSpec.Containers[0].Env = append(deploy.Spec.Template.PodSpec.Containers[0].Env, env)
			}
		}
		for _, service := range services {
			kubesServices = append(kubesServices, service)
		}
		deployments = append(deployments, deploy)
	}

	return kubernetes.Kubernetes{
		Services:    kubesServices,
		Deployments: deployments,
	}
}

func (s *ServiceYmlTransformer) ToServiceYml() service_yml.ServiceYml {
	return s.Base
}

func (s *ServiceYmlTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return docker_compose.DockerCompose{}
}
