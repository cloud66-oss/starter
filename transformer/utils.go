package transformer

import (
	"bufio"
	"unicode"
	"strings"
	"fmt"
	"os"
	"strconv"
	"math"
)

func finalFormat(lines []string) string {

	text := ""
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "cpu:") {
			lines[i] = formatCpu(lines[i])
		}
		text += lines[i] + "\n"
	}
	return text
}

func formatCpu(cpuString string) string {
	var i, auxInt, p int
	var auxString string
	//common.PrintlnTitle(cpuString)
	for i = 0; i < len(cpuString); i++ {
		if cpuString[i] == '\'' || cpuString[i] == '"' {
			p = i
			break
		}
	}
	for i = i + 1; i < len(cpuString); i++ {
		if cpuString[i] == '\'' || cpuString[i] == '"' {
			break
		} else {
			auxString += string(cpuString[i])
		}
	}
	auxFloat, err := strconv.ParseFloat(auxString, 64)
	checkError(err)
	fract := auxFloat - math.Floor(auxFloat)
	if auxFloat < 1 {
		auxInt = 1
	} else if fract < 0.5 {
		auxInt = int(math.Floor(auxFloat))
	} else {
		auxInt = int(math.Ceil(auxFloat))
	}
	cpuString = cpuString[:p] + strconv.Itoa(auxInt)
	return cpuString
}

func isEnv(line string) bool {
	for i := 0; i < len(line); i++ {
		if !unicode.IsSpace(rune(line[i])) {
			if line[i] == '-' {
				return true
			}
			break
		}
	}
	return false
}

func formatEnv_Vars(env string) string {
	var j int
	for j = 0; j < len(env); j++ {
		if env[j] == '-' {
			env = env[:j] + " " + env[j+1:]
		}
	}
	for ; j < len(env); j++ {
		if env[j] == '\'' {
			env = env[:j] + " " + env[j+1:]
		}
		if !unicode.IsSpace(rune(env[j])) {
			break
		}
	}
	for j = 0; j < len(env); j++ {
		if j+1 < len(env) {
			if env[j] == ':' {
				env = env[:j+1] + " " + env[j+1:]
				break
			}
			if env[j] == '=' {
				env = env[:j] + ": " + env[j+1:]
				break
			}
		}

	}
	for j = len(env) - 1; j >= 0; j-- {
		if env[j] == '-' {
			env = env[:j] + " " + env[j+1:]
		}
	}
	if strings.Contains(env, "\"\"") {
		return ""
	}
	return env
}

func readEnv_file(path string) map[string]string {
	var lines []string
	var env_vars map[string]string
	var key, value string
	envFile, err := os.Open(path)
	checkError(err)
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

func checkDB(image string) (string, bool) {
	db_check := []string{"mysql", "postgresql", "redis", "mongodb", "elasticsearch", "glusterfs", "influxdb", "rabbitmq", "sqlite"}
	var prefix string
	if len(image) < 5 {
		prefix = image
	} else {
		for i := 0; i < 4; i++ {
			prefix += string(image[i])
		}
	}
	for i := 0; i < len(db_check); i++ {
		if strings.Contains(image, db_check[i]) || strings.Contains(db_check[i], image) || strings.Contains(prefix, db_check[i]) || strings.Contains(db_check[i], prefix) {
			return db_check[i], true
		}
	}
	return "", false
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
