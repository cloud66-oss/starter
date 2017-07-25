package transform

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"unicode"

	"github.com/cloud66/starter/definitions/docker-compose"
	"github.com/cloud66/starter/definitions/service-yml"
	"github.com/cloud66/starter/common"
	"gopkg.in/yaml.v2"
	"strconv"
)

func readEnv_file(path string) map[string]string {
	var lines []string
	var env_vars map[string]string
	var key, value string
	envFile, err := os.Open(path)
	if err != nil {
		return env_vars
	}
	env_vars = make(map[string]string, 1)
	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < len(lines); i++ {

		if !isCommentLine(lines[i]) {

			key, value = getKeyValue(lines[i])
			env_vars[key] = value
		}
	}
	envFile.Close()
	return env_vars
}

func getKeyValue(line string) (string, string) {
	var key, value string
	var k int
	for k = 0; k < len(line); k++ {
		if !unicode.IsSpace(rune(line[k])) && line[k] != '"' {
			break
		}
	}
	for ; k < len(line); k++ {
		if line[k] == '=' || line[k] == '"' {
			break
		} else {
			key = string(append([]byte(key), line[k]))
		}
	}
	if line[k+1] == '=' && line[k+2] == '"' {
		k = k + 2
	} else if (line[k+1] == '=' && line[k+2] != '"') || (line[k+1] == '"') {
		k = k + 1
	}
	for k = k + 1; k < len(line); k++ {
		if line[k] == '\n' || line[k] == '"' {
			break
		} else {
			value = string(append([]byte(value), line[k]))
		}
	}

	return key, value
}

func isCommentLine(line string) bool {
	var i int
	for i = 0; i < len(line); i++ {
		if !unicode.IsSpace(rune(line[i])) {
			break
		}
	}
	if line != "" {
		if line[i] == '#' {
			return true
		}
	}
	return false
}

func dockerToServicePorts(exposed []int, dockerPorts docker_compose.Ports, shouldPrompt bool) service_yml.Ports {
	var servicePorts service_yml.Ports

	for _, expose := range exposed {
		servicePorts = append(servicePorts, service_yml.Port{
			Container: expose,
		})
	}
	for _, port := range dockerPorts {
		var servicePort service_yml.Port
		servicePort.Container = port.Target
		if port.Protocol == "tcp" {
			if shouldPrompt == true {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("\nYou have chosen a TCP protocol for the port published at %s - should it be mapped as HTTP, HTTPS or TCP ? : ", port.Published)
				var answer string
				answer, _ = reader.ReadString('\n')
				answer = strings.ToUpper(answer)
				if answer == "TCP\n" {
					servicePort.Tcp = port.Published
				}
				if answer == "HTTP\n" {
					servicePort.Http = port.Published
				}
				if answer == "HTTPS\n" {
					servicePort.Https = port.Published
				}
			} else {
				servicePort.Http = port.Published
			}
		} else {
			servicePort.Udp = port.Published
		}
		servicePorts = append(servicePorts, servicePort)
	}
	return servicePorts
}

func dockerToServiceVolumes(dockerVolumes docker_compose.Volumes) []string {
	var serviceVolumes []string

	for _, volume := range dockerVolumes {
		var temp string
		if volume.Type == "volume" {
			temp = volume.Source + ":" + volume.Target
			if volume.ReadOnly == true {
				temp = temp + ":ro"
			}
			if temp[0] != '/' && temp[0] != '$' {
				common.PrintlnWarning("Service.yml format does only support absolute path for volumes. Please modify for \"%s\"", temp)
				temp = "/" + temp
			}
		}
		serviceVolumes = append(serviceVolumes, temp)
	}
	return serviceVolumes
}

func dockerToServiceEnvVarFormat(service service_yml.Service) service_yml.Service {

	str, err := yaml.Marshal(service)
	service_yml.CheckError(err)
	for i := 0; i < len(str); i++ {
		if str[i] == '{' && str[i-1] == '$' {
			str = []byte(string(str[:i-1]) + "_env(" + string(str[i+1:]))
			for ; i < len(str); i++ {
				if str[i] == '}' {
					str[i] = ')'
					break
				}
			}
		}
	}
	var newService service_yml.Service
	err = yaml.Unmarshal(str, &newService)
	service_yml.CheckError(err)

	return newService
}

func dockerToServiceStopGrace(str string) int {
	if str != "" {
		var stopInt int
		var err error
		if !unicode.IsDigit(rune(str[len(str)-1])) {
			stopInt, err = strconv.Atoi(str[:len(str)-1])
			if err != nil{
				stopInt = 30 //usual used number in case the user needs a stop grace - can be modified afterwards
			}
			return stopInt
		}else {
			stopInt, err = strconv.Atoi(str)
			if err!=nil{
				stopInt = 30 //usual used number in case the user needs a stop grace - can be modified afterwards
			}
			return stopInt
		}
	}
	return 0
}

func getDockerToServiceWarnings(service docker_compose.Service) {
	if service.CapAdd != nil {
		common.PrintlnWarning("Service.yml format does not support \"cap_add\" at the moment")
	}
	if service.CapDrop != nil {
		common.PrintlnWarning("Service.yml format does not support \"cap_drop\" at the moment")
	}
	if service.ContainerName != "" {
		common.PrintlnWarning("Service.yml format does not support \"container_name\" at the moment")
	}
	if service.CgroupParent != "" {
		common.PrintlnWarning("Service.yml format does not support \"cgroup_parent\" at the moment")
	}
	if service.Devices != nil {
		common.PrintlnWarning("Service.yml format does not support \"devices\" at the moment")
	}
	if service.Links != nil {
		common.PrintlnWarning("Service.yml format does not support \"links\" at the moment")
	}
	if service.Dns != nil {
		common.PrintlnWarning("Service.yml format does not support \"dns\" at the moment")
	}
	if service.DnsSearch != nil {
		common.PrintlnWarning("Service.yml format does not support \"dns_search\" at the moment")
	}
	if service.ExtraHosts != nil {
		common.PrintlnWarning("Service.yml format does not support \"hosts\" at the moment")
	}
	if service.Isolation != "" {
		common.PrintlnWarning("Service.yml format does not support \"isolation\" at the moment")
	}
	if service.Networks.Aliases != nil {
		common.PrintlnWarning("Service.yml format does not support \"networks\" at the moment")
	}
	if service.Secrets != nil {
		common.PrintlnWarning("Service.yml format does not support \"secrets\" at the moment")
	}
	if service.SecurityOpt != nil {
		common.PrintlnWarning("Service.yml format does not support \"security_opt\" at the moment")
	}
	if service.UsernsMode != "" {
		common.PrintlnWarning("Service.yml format does not support \"userns_mode\" at the moment")
	}
	if service.Ulimits.Nproc.Soft != 0 || service.Ulimits.Nproc.Hard != 0 || service.Ulimits.Nofile.Soft != 0 || service.Ulimits.Nofile.Hard != 0 {
		common.PrintlnWarning("Service.yml format does not support \"ulimits\" at the moment")
	}
	if service.Healthcheck.Interval != "" || service.Healthcheck.Test != nil || service.Healthcheck.Timeout != "" || service.Healthcheck.Disable == true {
		common.PrintlnWarning("Service.yml format does not support \"healthcheck\" at the moment")
	}
	if service.Logging.Driver != "" || service.Logging.Options != nil {
		common.PrintlnWarning("Service.yml format does not support \"logging\" at the moment")
	}
	if service.Deploy.Resources.Limits.Cpus != "" || service.Deploy.Resources.Limits.Memory != "" {
		common.PrintlnWarning("Service.yml format does not support \"resources limitations and reservations\" for deploy at the moment, try using \"cpu_shares\" and \"mem_limit\" instead. ")
	}
	if service.Deploy.UpdateConfig.Delay != "" || service.Deploy.UpdateConfig.Parallelism != 0 {
		common.PrintlnWarning("Service.yml format does not support \"update_config\" at the moment")
	}
	if service.Deploy.Placement.Constraints != nil {
		common.PrintlnWarning("Service.yml format does not support \"placement constraints\" at the moment")
	}
}
