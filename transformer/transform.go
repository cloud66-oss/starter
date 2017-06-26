package transformer

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"

	"github.com/cloud66/starter/common"
)

//main transformation format function
func Transformer(filename string, formatTarget string) error {

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

		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(d)

	} else {
		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(dockerCompose.Services)
	}

	file, err := yaml.Marshal(serviceYaml)

	file = []byte("# Generated with <3 by Cloud66\n\n"+string(file))

	err = ioutil.WriteFile("service.yml", file, 0644)
	if err != nil { //tvbot
		log.Fatalf("ioutil.WriteFile: %v", err)
	}

	return err

}

func copyToServiceYML(d map[string]DockerService) (map[string]ServiceYMLService, []string) {

	serviceYaml := ServiceYml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      make([]string, 0),
	}
	var isDB bool
	var err error
	var dbs []string

	for k, v := range d {
		var current_db string
		isDB = false

		var gitURL, gitBranch, buildRoot string

		if v.Image != "" {
			current_db, isDB = checkDB(v.Image)

		} else {

			var gitPath string
			gitPath, err = common.GitRootDir("/")
			if err != nil {

			}
			hasGit := common.HasGit(gitPath)

			if hasGit {
				gitURL = common.RemoteGitUrl(gitPath)
				gitBranch = common.LocalGitBranch(gitPath)
				buildRoot, err = common.PathRelativeToGitRoot(gitPath)
			}
		}
		if isDB {
			dbs = append(dbs, current_db)
		}
		if !isDB {


			if v.Deploy.Resources.Limits.Cpus!="" || v.Deploy.Resources.Limits.Memory!="" || v.Deploy.Resources.Reservations.Cpus!="" ||v.Deploy.Resources.Reservations.Memory!=""{
				common.PrintlnWarning("Service.yml format does not support \"resources limitations and reservations\" for deploy at the moment, try using \"cpu_shares\" and \"mem_limit\" instead. ")

			}


			var serviceYamlService ServiceYMLService
			serviceYamlService.GitRepo = gitURL
			serviceYamlService.GitBranch = gitBranch
			serviceYamlService.BuildRoot = buildRoot
			serviceYamlService.BuildCommand = v.BuildCommand.BuildCommand
			serviceYamlService.Command = v.Command.Command
			serviceYamlService.Image = v.Image
			serviceYamlService.Requires = v.Depends_on
			serviceYamlService.Volumes = handleVolumes(v.Volumes.Volumes, v.Volumes.LongSyntax)
			serviceYamlService.Ports = handlePorts(v.Expose, v.Ports.Port, v.Ports.ShortSyntax)
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


			serviceYamlService.Tags = make(map[string]string,1)
			for key, w := range v.Deploy.Labels {
				serviceYamlService.Tags[key] = w
			}

			if v.EnvFile.EnvFile != nil {
				var lines map[string]string
				for i := 0; i < len(v.EnvFile.EnvFile); i++ {
					lines = readEnv_file(v.EnvFile.EnvFile[i])
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
