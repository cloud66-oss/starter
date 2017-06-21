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

	dockerCompose := Docker_compose{
		Services: make(map[string]docker_Service),
		Version:  "",
	}

	serviceYaml := Serviceyml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      make([]string, 0),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &dockerCompose); err != nil {
		fmt.Println(err.Error())
	}

	if len(dockerCompose.Services) == 0 {
		d := make(map[string]docker_Service)
		err = yaml.Unmarshal([]byte(yamlFile), &d)
		CheckError(err)

		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(d)

	} else {

		serviceYaml.Services, serviceYaml.Dbs = copyToServiceYML(dockerCompose.Services)
	}
	if len(serviceYaml.Dbs) != 0 {
		if serviceYaml.Dbs[len(serviceYaml.Dbs)-1] == "" {
			serviceYaml.Dbs = serviceYaml.Dbs[:len(serviceYaml.Dbs)-1]
		}
	}

	file, err := yaml.Marshal(serviceYaml)

	err = ioutil.WriteFile("service.yml", file, 0644)
	if err != nil {
		log.Fatalf("ioutil.WriteFile: %v", err)
	}

	service_yml, er := os.OpenFile(formatTarget, os.O_RDWR, 0644)
	CheckError(er)

	// write some text to file
	_, err = service_yml.WriteString(string(file))
	CheckError(err)

	// save changes
	err = service_yml.Sync()
	CheckError(err)

	service_yml.Close()

	//format long syntax ports
	service_yml, _ = os.Open(formatTarget)
	defer service_yml.Close()

	var lines []string
	scanner := bufio.NewScanner(service_yml)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//final format for ENV_VARS, CPU and PORTS
	text := finalFormat(lines)

	//write the final service.yml
	service_yml, _ = os.Create(formatTarget)
	service_yml, er = os.OpenFile(formatTarget, os.O_RDWR, 0644)

	_, err = service_yml.WriteString(text)

	CheckError(err)

	return err

}

func copyToServiceYML(d map[string]docker_Service) (map[string]ServiceYMLService, []string) {

	serviceYaml := Serviceyml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      make([]string, 0),
	}
	var isDB bool
	var err error
	var dbs []string

	for k, v := range d {
		var current_db string
		isDB = false

		if v.Image != "" {
			current_db, isDB = checkDB(v.Image)
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
				longSyntaxPorts = []string{}
				for i := 0; i < len(v.Ports.Port); i++ {

					longSyntax := ""
					longSyntax = "target: " + v.Ports.Port[i].Target + "\n"

					if v.Ports.Port[i].Protocol == "udp" {
						longSyntax += "udp: " + v.Ports.Port[i].Published
						longSyntaxPorts = append(longSyntaxPorts, longSyntax)
					} else if v.Ports.Port[i].Protocol == "tcp" {
						reader := bufio.NewReader(os.Stdin)
						fmt.Printf("\nYou have chosen a TCP protocol for the port published at %s - should it be mapped as HTTP, HTTPS or TCP ? : ", v.Ports.Port[i].Published)
						var answer string //tvbot
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

			var gitPath string
			gitPath, err = common.GitRootDir("/")
			if err != nil {

			}
			hasGit := common.HasGit(gitPath)

			var gitURL, gitBranch, buildRoot string
			if hasGit {
				gitURL = common.RemoteGitUrl(gitPath)
				gitBranch = common.LocalGitBranch(gitPath)
				buildRoot, err = common.PathRelativeToGitRoot(gitPath)
			}

			var serviceYamlService ServiceYMLService
			serviceYamlService.GitRepo = gitURL
			serviceYamlService.GitBranch = gitBranch
			serviceYamlService.BuildRoot = buildRoot
			serviceYamlService.BuildCommand = v.Build_Command.Build_Command
			serviceYamlService.Command = v.Command.Command
			serviceYamlService.Image = v.Image
			serviceYamlService.Requires = v.Depends_on
			serviceYamlService.Volumes = v.Volumes.Volumes
			serviceYamlService.StopGrace = v.Stop_grace_period
			serviceYamlService.WorkDir = v.Working_dir
			serviceYamlService.EnvVars = v.EnvVars
			serviceYamlService.Tags = v.Labels
			serviceYamlService.DockerfilePath = v.Build_Command.Build.Dockerfile
			serviceYamlService.Privileged = v.Privileged
			serviceYamlService.Constraints = Constraints{
				Resources: Resources{
					Memory: v.MemLimit,
					Cpu:    v.CpuShares,
				},
			}
			serviceYamlService.Ports = longSyntaxPorts
			for key,w := range v.Deploy.Labels{
				serviceYamlService.Tags[key]=w
			}

			if v.Env_file.Env_file != nil {
				var lines map[string]string
				for i := 0; i < len(v.Env_file.Env_file); i++ {
					lines = readEnv_file(v.Env_file.Env_file[i])
					for j,w := range lines{
						if j!=""{
							serviceYamlService.EnvVars[j]=w
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
