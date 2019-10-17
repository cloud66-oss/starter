package transform

import (
	"bufio"
	"os"
	"unicode"

	"github.com/cloud66-oss/starter/definitions/service-yml"
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
	if k < len(line)-2 {
		if line[k+1] == '=' && line[k+2] == '"' {
			k = k + 2
		} else if (line[k+1] == '=' && line[k+2] != '"') || (line[k+1] == '"') {
			k = k + 1
		}
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
			if err != nil {
				stopInt = 30 //usual used number in case the user needs a stop grace - can be modified afterwards
			}
			return stopInt
		} else {
			stopInt, err = strconv.Atoi(str)
			if err != nil {
				stopInt = 30 //usual used number in case the user needs a stop grace - can be modified afterwards
			}
			return stopInt
		}
	}
	return 0
}
