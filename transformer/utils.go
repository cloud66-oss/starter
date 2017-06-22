package transformer

import (
	"bufio"
	"unicode"
	"strings"
	"fmt"
	"os"
)

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
		if strings.Contains(image, db_check[i]) || strings.Contains(db_check[i], image) {
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
		os.Exit(0)
	}
}
