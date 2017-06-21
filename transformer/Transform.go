package transformer

import (
	"bufio"
	"strings"
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

	//reformat the output from long syntax ports
	text := string(file) //tvbot

	lines := strings.Split(text, "\n")
	file = finalFormat(lines)

	err = ioutil.WriteFile("service.yml", file, 0644)
	if err != nil {
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
			var longSyntaxPorts []string
			longSyntaxPorts = v.Expose
			if len(v.Ports.ShortSyntax) > 0 {
				for i := 0; i < len(v.Ports.ShortSyntax); i++ {
					longSyntaxPorts = append(longSyntaxPorts, v.Ports.ShortSyntax[i])
				}
			} else {
				for i := 0; i < len(v.Ports.Port); i++ {

					longSyntax := ""
					longSyntax = "target: " + v.Ports.Port[i].Target + "\n"

					if v.Ports.Port[i].Protocol == "udp" {
						longSyntax += "udp: " + v.Ports.Port[i].Published
						longSyntaxPorts = append(longSyntaxPorts, longSyntax)
					} else if v.Ports.Port[i].Protocol == "tcp" {
						reader := bufio.NewReader(os.Stdin)
						fmt.Printf("\nYou have chosen a TCP protocol for the port published at %s - should it be mapped as HTTP, HTTPS or TCP ? : ", v.Ports.Port[i].Published)
						var answer string
						answer, _ = reader.ReadString('\n')
						answer = strings.ToUpper(answer)
						if answer == "TCP\n" {
							longSyntax += "tcp: " + v.Ports.Port[i].Published
						} else if answer == "HTTP\n" {
							longSyntax += "http: " + v.Ports.Port[i].Published
						} else if answer == "HTTPS\n" {
							longSyntax += "http: " + v.Ports.Port[i].Published
						}
						longSyntaxPorts = append(longSyntaxPorts, longSyntax)
					}

				}
			}

			var serviceYamlService ServiceYMLService
			serviceYamlService.GitRepo = gitURL
			serviceYamlService.GitBranch = gitBranch
			serviceYamlService.BuildRoot = buildRoot
			serviceYamlService.BuildCommand = v.BuildCommand.BuildCommand
			serviceYamlService.Command = v.Command.Command
			serviceYamlService.Image = v.Image
			serviceYamlService.Requires = v.Depends_on
			serviceYamlService.Volumes = v.Volumes.Volumes
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
			serviceYamlService.Ports = longSyntaxPorts
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
