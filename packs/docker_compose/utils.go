package docker_compose

import (
	"bufio"
	"unicode"
	"strings"
	"fmt"
	"os"
)

func handleEnvVarsFormat(text []byte) string {

	for i := 0; i < len(text)-1; i++ {
		if text[i] == '$' && text[i+1] == '{' {
			text = []byte(string(text[:i]) + "_env(" + string(text[i+2:]))
			for ; i < len(text); i++ {
				if text[i] == '-' {
					if text[i+1] == '.' && text[i+2] == '}' {
						text = []byte(string(text[:i+1]) + "Default" + string(text[i+2:]))
					}
					if text[i-1] == ':' {
						text = []byte(string(text[:i]) + string(text[i+1:]))
					} else {
						text[i] = ':'
					}
				}
				if text[i] == '}' {
					text[i] = ')'
					break
				}
			}
		}
	}
	return string(text)
}

func handleVolumes(shortSyntax []string, longSyntax []LongSyntaxVolume) []interface{} {

	var longSyntaxVolumes []interface{}

	if len(shortSyntax) > 0 {
		for i := 0; i < len(shortSyntax); i++ {
			if shortSyntax[i][0] == '.'{
				shortSyntax[i]= shortSyntax[i][1:]
			}
			if shortSyntax[i][0]!='/'{
				shortSyntax[i] = "/"+shortSyntax[i]
			}
			longSyntaxVolumes = append(longSyntaxVolumes, shortSyntax[i])
		}
	}

	var tempString string
	for i := 0; i < len(longSyntax); i++ {
		if longSyntax[i].Type == "volume" {
			if longSyntax[i].ReadOnly == true {
				tempString = longSyntax[i].Source + ":" + longSyntax[i].Target + ":ro"
			} else {
				tempString = longSyntax[i].Source + ":" + longSyntax[i].Target
			}
		}
		if tempString[0] == '.'{
			tempString = tempString[1:]
		}
		if tempString[0]!='/'{
			tempString = "/"+tempString
		}
		longSyntaxVolumes = append(longSyntaxVolumes, tempString)
	}

	return longSyntaxVolumes
}

func formatShortPorts(dockerPort string) string {
	var servicePort, aux, aux2 string
	var i int

	for i = 0; i < len(dockerPort); i++ {
		if unicode.IsDigit(rune(dockerPort[i])) {
			break
		}
	}
	for ; i < len(dockerPort); i++ {
		if dockerPort[i] == ':' || dockerPort[i] == '\n' {
			break
		} else {
			aux = string(append([]byte(aux), dockerPort[i]))
		}
	}
	for i = i + 1; i < len(dockerPort); i++ {
		if dockerPort[i] == ':' || dockerPort[i] == '\n' {
			break
		} else {
			aux2 = string(append([]byte(aux2), dockerPort[i]))
		}
	}
	if i < len(dockerPort)-1 {
		servicePort = aux2 + ":" + aux + dockerPort[i:]
	} else {
		servicePort = aux2 + ":" + aux
	}
	if servicePort[0] == ':' {
		servicePort = servicePort[1:]
	}
	return servicePort
}

func handlePorts(expose []string, longSyntax []Port, shortSyntax []string, shouldPrompt bool) []interface{} {

	var longSyntaxPorts []interface{}

	for i := 0; i < len(expose); i++ {
		longSyntaxPorts = append(longSyntaxPorts, expose[i])
	}
	if len(shortSyntax) > 0 {
		for i := 0; i < len(shortSyntax); i++ {
			shortSyntax[i] = formatShortPorts(shortSyntax[i])
			longSyntaxPorts = append(longSyntaxPorts, shortSyntax[i])
		}
	}

	for i := 0; i < len(longSyntax); i++ {

		var serviceyml_longsyntax ServicePort
		serviceyml_longsyntax.Container = longSyntax[i].Target

		if longSyntax[i].Protocol == "tcp" {
			if shouldPrompt == true {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("\nYou have chosen a TCP protocol for the port published at %s - should it be mapped as HTTP, HTTPS or TCP ? : ", longSyntax[i].Published)
				var answer string
				answer, _ = reader.ReadString('\n')
				answer = strings.ToUpper(answer)
				if answer == "TCP\n" {
					serviceyml_longsyntax.Tcp = longSyntax[i].Published
				}
				if answer == "HTTP\n" {
					serviceyml_longsyntax.Http = longSyntax[i].Published
				}
				if answer == "HTTPS\n" {
					serviceyml_longsyntax.Https = longSyntax[i].Published
				}
			} else {
				serviceyml_longsyntax.Http = longSyntax[i].Published
			}
		} else {
			serviceyml_longsyntax.Udp = longSyntax[i].Published
		}
		longSyntaxPorts = append(longSyntaxPorts, serviceyml_longsyntax)
	}

	return longSyntaxPorts
}

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
		if !unicode.IsSpace(rune(line[k])) {
			break
		}
	}
	for ; k < len(line); k++ {
		if line[k] == '=' {
			break
		} else {
			key = string(append([]byte(key), line[k]))
		}
	}
	for k = k + 1; k < len(line); k++ {
		if line[k] == '\n' {
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

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
