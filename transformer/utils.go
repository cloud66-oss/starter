package transformer

import (
	"bufio"
	"unicode"
	"strings"
	"fmt"
	"os"
)

func handleVolumes(shortSyntax []string, longSyntax []LongSyntaxVolume) []interface{} {

	var longSyntaxVolumes []interface{}

	if len(shortSyntax) > 0 {
		for i := 0; i < len(shortSyntax); i++ {
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
		longSyntaxVolumes = append(longSyntaxVolumes, tempString)
	}

	return longSyntaxVolumes
}

func handlePorts(expose []string, longSyntax []Port, shortSyntax []string) []interface{}{

	var longSyntaxPorts []interface{}

	for i := 0; i < len(expose); i++ {
		longSyntaxPorts = append(longSyntaxPorts, expose[i])
	}
	if len(shortSyntax) > 0 {
		for i := 0; i < len(shortSyntax); i++ {
			longSyntaxPorts = append(longSyntaxPorts, shortSyntax[i])
		}
	}

	for i := 0; i < len(longSyntax); i++ {

		var serviceyml_longsyntax ServicePort
		serviceyml_longsyntax.Container = longSyntax[i].Target

		if longSyntax[i].Protocol == "tcp" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("\nYou have chosen a TCP protocol for the port published at %s - should it be mapped as HTTP, HTTPS or TCP ? : ", longSyntax[i].Published)
			var answer string
			answer, _ = reader.ReadString('\n')
			answer = strings.ToUpper(answer)
			if answer == "TCP\n"{
				serviceyml_longsyntax.Tcp = longSyntax[i].Published
			}
			if answer == "HTTP\n"{
				serviceyml_longsyntax.Http = longSyntax[i].Published
			}
			if answer == "HTTPS\n"{
				serviceyml_longsyntax.Https = longSyntax[i].Published
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
	CheckError(err)
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
	db_check := []string{"mysql", "postgresql", "redis", "mongodb", "elasticsearch", "glusterfs", "influxdb", "rabbitmq", "sqlite", "postgres", "mongo", "influx"}

	for i := 0; i < len(db_check); i++ {
		if strings.Contains(image, db_check[i]) {
			switch db_check[i] {
			case "postgres":
				return "postgresql", true
			case "mongo":
				return "mongodb", true
			case "influx":
				return "influxdb", true
			default:
				return db_check[i], true
			}
		}
	}
	return "", false
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
