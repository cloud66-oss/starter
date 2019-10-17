package packs

import (
	"github.com/cloud66-oss/starter/common"
	"github.com/cloud66-oss/starter/definitions/service-yml"
	"strconv"
)

type ServiceYAMLContextBase struct {
	Services []*common.Service
	Dbs      []common.Database
}

func (s *ServiceYAMLContextBase) GenerateFromServiceYml(serviceYml service_yml.ServiceYml) error {
	s.Services = make([]*common.Service, 0)
	s.Dbs = make([]common.Database, 0)

	for name, service := range serviceYml.Services {
		//migrate from service_yml.Service to common.Service

		ports := make([]*common.PortMapping, 0)
		// create for to populate ports

		for _, port := range service.Ports {
			ports = append(ports, &common.PortMapping{
				Container: strconv.Itoa(port.Container),
				HTTP:      strconv.Itoa(port.Http),
				HTTPS:     strconv.Itoa(port.Https),
				TCP:       strconv.Itoa(port.Tcp),
				UDP:       strconv.Itoa(port.Udp),
			})
		}

		envs := make([]*common.EnvVar, 0)

		for key, value := range service.EnvVars {
			envs = append(envs, &common.EnvVar{
				Key:   key,
				Value: value,
			})
		}

		s.Services = append(s.Services, &common.Service{
			Name:          name,
			GitRepo:       service.GitUrl,
			GitBranch:     service.GitBranch,
			Command:       service.Command,
			BuildCommand:  service.BuildCommand,
			DeployCommand: service.DeployCommand,
			Ports:         ports,
			EnvVars:       envs,
			BuildRoot:     service.BuildRoot,
			Databases:     make([]common.Database, 0),
		})
	}

	for _, database := range serviceYml.Databases {
		//migrate from string to common.Database
		s.Dbs = append(s.Dbs, common.Database{
			Name:        database,
			DockerImage: database,
		})
	}

	return nil
}
