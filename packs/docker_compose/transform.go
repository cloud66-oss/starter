package docker_compose

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"github.com/cloud66/starter/common"
)

//main transformation format function
func Transformer(filename string, formatTarget string, gitURL string, gitBranch string, shouldPrompt bool) error {

	var err error
	_, err = os.Stat(formatTarget)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(formatTarget)
		CheckError(err)
		defer file.Close()
	} else {
		common.PrintError("File %s already exists. Will be overwritten.\n", formatTarget)
	}

	yamlFile, err := ioutil.ReadFile(filename)

	dockerCompose := DockerCompose{
		Services: make(map[string]DockerService),
		Version:  "",
	}

	serviceYaml := ServiceYml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      make([]string, 0),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &dockerCompose); err != nil {
		fmt.Println(err.Error())
	}

	// Due to the fact that docker-compose versions prior to v2.0 did not
	// require a structured yml file with the top-level as being occupied
	// by "services:", the previous unmarshal might not work and the code
	// checks whether this is the case and if so, tries to unmarshal w/o
	// the top-level

	if len(dockerCompose.Services) == 0 {
		d := make(map[string]DockerService)
		err = yaml.Unmarshal([]byte(yamlFile), &d)
		CheckError(err)

		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(d, gitURL, gitBranch, shouldPrompt, filename)

	} else {
		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(dockerCompose.Services, gitURL, gitBranch, shouldPrompt, filename)
	}

	file, err := yaml.Marshal(serviceYaml)

	file = []byte("# Generated with <3 by Cloud66\n\n" + handleEnvVarsFormat(file))

	err = ioutil.WriteFile(formatTarget, file, 0644)
	if err != nil {
		return err
	}

	return nil

}

func copyToServiceYML(d map[string]DockerService, gitURL string, gitBranch string, shouldPrompt bool, filepath string) (map[string]ServiceYMLService, []string) {

	serviceYaml := ServiceYml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      make([]string, 0),
	}
	var isDB bool
	var err error
	var dbs []string

	var dbServicesNames []string
	dbServicesNames = make([]string, 1)
	servicesNames := make([]string, 1)

	for k, v := range d {

		isDB = false

		if v.Image != "" {
			_, isDB = checkDB(v.Image)
		}
		if isDB {
			dbServicesNames = append(dbServicesNames, k)
		}
		servicesNames = append(servicesNames, k)
	}

	for k, v := range d {
		var current_db string
		isDB = false

		//var gitURL, gitBranch string
		var buildRoot string

		if v.Image != "" {
			current_db, isDB = checkDB(v.Image)

		} else {

			var gitPath string
			gitPath, err = common.GitRootDir("/")
			if err != nil {

			}
			hasGit := common.HasGit(gitPath)

			if hasGit {
				//gitURL = common.RemoteGitUrl(gitPath)
				//gitBranch = common.LocalGitBranch(gitPath)
				buildRoot, err = common.PathRelativeToGitRoot(gitPath)
			}
		}
		if isDB {
			dbServicesNames = append(dbServicesNames, k)
			dbs = append(dbs, current_db)
		}
		if !isDB {

			if v.Deploy.Resources.Limits.Cpus != "" || v.Deploy.Resources.Limits.Memory != "" || v.Deploy.Resources.Reservations.Cpus != "" || v.Deploy.Resources.Reservations.Memory != "" {
				common.PrintlnWarning("Service.yml format does not support \"resources limitations and reservations\" for deploy at the moment, try using \"cpu_shares\" and \"mem_limit\" instead. ")

			}

			var serviceYamlService ServiceYMLService

			//Set Git stuff and BuildRoot
			serviceYamlService.GitRepo = gitURL
			serviceYamlService.GitBranch = gitBranch
			if v.BuildCommand.BuildRoot != "" {
				serviceYamlService.BuildRoot = v.BuildCommand.BuildRoot
			} else if v.BuildCommand.Build.Context != "" {
				serviceYamlService.BuildRoot = v.BuildCommand.Build.Context
			} else {
				serviceYamlService.BuildRoot = buildRoot
			}

			//Do not let a service require a db as service.yml runs dbs first anyway
			for i := 0; i < len(v.Links.Links); i++ {
				v.Depends_on = append(v.Depends_on, v.Links.Links[i])
			}
			for i := 0; i < len(v.Depends_on); i++ {
				for j := 0; j < len(dbServicesNames); j++ {
					if i < len(v.Depends_on) {
						if v.Depends_on[i] == dbServicesNames[j] {
							v.Depends_on = append(v.Depends_on[:i], v.Depends_on[i+1:]...)
						}
					}
				}
			}

			serviceYamlService.Command = v.Command.Command
			serviceYamlService.Image = v.Image
			serviceYamlService.Requires = v.Depends_on
			serviceYamlService.Volumes = handleVolumes(v.Volumes.Volumes, v.Volumes.LongSyntax)
			serviceYamlService.Ports = handlePorts(v.Expose, v.Ports.Port, v.Ports.ShortSyntax, shouldPrompt)
			serviceYamlService.StopGrace = v.Stop_grace_period
			serviceYamlService.WorkDir = v.Working_dir
			serviceYamlService.EnvVars = v.EnvVars.EnvVars
			serviceYamlService.Tags = v.Labels
			serviceYamlService.DockerfilePath = v.BuildCommand.Build.Dockerfile
			serviceYamlService.Privileged = v.Privileged
			serviceYamlService.Constraints = Constraints{
				Resources: Resources{
					Memory: v.MemLimit,
					Cpu:    v.CpuShares,
				},
			}

			serviceYamlService.Tags = make(map[string]string, 1)
			for key, w := range v.Deploy.Labels {
				serviceYamlService.Tags[key] = w
			}

			if v.EnvFile.EnvFile != nil {
				var lines map[string]string
				for i := 0; i < len(v.EnvFile.EnvFile); i++ {
					path := filepath[0:len(filepath)-len("docker-compose.yml")]
					lines = readEnv_file(path + v.EnvFile.EnvFile[i])
					for j, w := range lines {
						if j != "" {
							serviceYamlService.EnvVars[j] = w
						}
					}
				}
			}

			if serviceYamlService.Image != "" {
				serviceYamlService.GitRepo = ""
				serviceYamlService.GitBranch = ""
				serviceYamlService.BuildRoot = ""
			}
			serviceYaml.Services[k] = serviceYamlService
		}
	}

	return serviceYaml.Services, dbs
}
