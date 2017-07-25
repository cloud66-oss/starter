package transform

import (
	"strings"

	"github.com/cloud66/starter/definitions/docker-compose"
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/service-yml"
	"github.com/cloud66/starter/common"
)

type DockerComposeTransformer struct {
	Base docker_compose.DockerCompose
}

func (d *DockerComposeTransformer) ToKubernetes() kubernetes.Kubernetes {
	return kubernetes.Kubernetes{}
}

func (d *DockerComposeTransformer) ToServiceYml(gitURL string, gitBranch string, shouldPrompt bool, filepath string) service_yml.ServiceYml {
	serviceYaml := service_yml.ServiceYml{
		Services:  make(map[string]service_yml.Service),
		Databases: make([]string, 0),
	}
	var err error

	for k, v := range d.Base.Services {
		var buildRoot string
		var gitPath string
		gitPath, err = common.GitRootDir("/")
		if err != nil {

		}
		hasGit := common.HasGit(gitPath)

		getDockerToServiceWarnings(v)

		if hasGit {
			buildRoot, err = common.PathRelativeToGitRoot(gitPath)
		}

		if v.Deploy.Resources.Limits.Cpus != "" || v.Deploy.Resources.Limits.Memory != "" || v.Deploy.Resources.Reservations.Cpus != "" || v.Deploy.Resources.Reservations.Memory != "" {
			common.PrintlnWarning("Service.yml format does not support \"resources limitations and reservations\" for deploy at the moment, try using \"cpu_shares\" and \"mem_limit\" instead. ")

		}

		var serviceYamlService service_yml.Service

		//Set Git stuff and BuildRoot
		serviceYamlService.GitUrl = gitURL
		serviceYamlService.GitBranch = gitBranch
		if v.Build.Context != "" {
			serviceYamlService.BuildRoot = v.Build.Context
		} else {
			serviceYamlService.BuildRoot = buildRoot
		}

		if serviceYamlService.BuildRoot == "." {
			serviceYamlService.BuildRoot = ""
		}

		serviceYamlService.Command = strings.Join(v.Command, " ")
		serviceYamlService.Image = v.Image
		serviceYamlService.Requires = v.Depends_on
		serviceYamlService.Volumes = dockerToServiceVolumes(v.Volumes)
		serviceYamlService.Ports = dockerToServicePorts(v.Expose, v.Ports, shouldPrompt)


		serviceYamlService.StopGrace = dockerToServiceStopGrace(v.Stop_grace_period)
		serviceYamlService.WorkDir = v.Working_dir
		serviceYamlService.EnvVars = v.Environment
		serviceYamlService.Tags = v.Labels
		serviceYamlService.DockerfilePath = v.Build.Dockerfile
		serviceYamlService.Privileged = v.Privileged
		serviceYamlService.Constraints = service_yml.Constraints{
			Resources: service_yml.Resources{
				Memory: v.MemLimit,
				Cpu:    v.CpuShares,
			},
		}

		serviceYamlService.Tags = make(map[string]string, 1)
		for key, w := range v.Deploy.Labels {
			serviceYamlService.Tags[key] = w
		}

		if len(v.EnvFile) > 0 {
			var lines map[string]string
			for i := 0; i < len(v.EnvFile); i++ {
				path := filepath[0:len(filepath)-len("docker-compose.yml")]
				lines = readEnv_file(path + v.EnvFile[i])
				for j, w := range lines {
					if j != "" {
						if len(serviceYamlService.EnvVars)==0{
							serviceYamlService.EnvVars = make(map[string]string,1)
						}
						serviceYamlService.EnvVars[j] = w
					}
				}
			}
		}

		if serviceYamlService.Image != "" {
			serviceYamlService.GitUrl = ""
			serviceYamlService.GitBranch = ""
			serviceYamlService.BuildRoot = ""
		}
		serviceYamlService = dockerToServiceEnvVarFormat(serviceYamlService)
		serviceYaml.Services[k] = serviceYamlService
	}

	return serviceYaml
}

func (d *DockerComposeTransformer) ToDockerCompose() docker_compose.DockerCompose {
	return d.Base
}
