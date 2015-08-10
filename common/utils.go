package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/mgutz/ansi"
)

var (
	procfileRegex = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")
	envVarPattern = "\\$([A-Z_]+[A-Z0-9_]*)"
	envVarRegex   = regexp.MustCompile(envVarPattern)

	MsgTitle string = ansi.ColorCode("green+h")
	MsgL0    string = ansi.ColorCode("magenta")
	MsgL1    string = ansi.ColorCode("white")
	MsgL2    string = ansi.ColorCode("black+h")
	MsgReset string = ansi.ColorCode("reset")
	MsgError string = ansi.ColorCode("red")
	MsgWarn  string = ansi.ColorCode("yellow")
)

type Process struct {
	Name    string
	Command string
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ContainsString(slice []string, desired string) bool {
	for _, item := range slice {
		if item == desired {
			return true
		}
	}
	return false
}

func CompareVersions(desired string, actual string) (bool, error) {
	act, err := version.NewVersion(actual)
	if err != nil {
		return false, err
	}

	des, err := version.NewConstraint(desired)
	if err != nil {
		return false, err
	}

	return des.Check(act), nil
}

func ParseProcfile(procfile string) ([]*Process, error) {
	procs := []*Process{}

	buf, err := ioutil.ReadFile(procfile)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(buf), "\n") {
		if matches := procfileRegex.FindStringSubmatch(line); matches != nil {
			name, command := matches[1], matches[2]
			procs = append(procs, &Process{Name: name, Command: command})
		}
	}

	return procs, nil
}

func ParseEnvironmentVariables(line string) (string, error) {
	line = envVarRegex.ReplaceAllString(line, "_env:$1")

	return line, nil
}

func ParseUniqueInt(line string) (string, error) {
	return strings.Replace(line, "{{UNIQUE_INT}}", "_unique:int", -1), nil
}

func ParsePort(command string) (hasFound bool, port string) {
	portRegexp := regexp.MustCompile(`(?:-p|--port=)[[:blank:]]*(\d+)`)
	ports := portRegexp.FindAllStringSubmatch(command, -1)
	if len(ports) != 1 {
		return false, ""
	} else {
		return true, ports[0][1]
	}
}

func RemovePortIfEnvVar(command string) string {
	portEnvVarRegexp := regexp.MustCompile(`[[:blank:]]*(-p|--port=)[[:blank:]]*` + envVarPattern)
	return portEnvVarRegexp.ReplaceAllString(command, "")
}

func AskUser(message string) string {
	answer := ""
	for strings.TrimSpace(answer) == "" {
		fmt.Print(MsgL1, fmt.Sprintf(" %s: ", message), MsgReset)
		fmt.Scanln(&answer)
	}
	return answer
}

func AskUserWithDefault(message string, defaultValue string, shouldPrompt bool) string {
	if !shouldPrompt {
		return defaultValue
	}

	printedDefaultValue := defaultValue
	if printedDefaultValue == "" {
		printedDefaultValue = "default: none"
	}

	fmt.Print(MsgL1, fmt.Sprintf(" %s [%s] ", message, printedDefaultValue), MsgReset)
	value := ""
	if _, err := fmt.Scanln(&value); err != nil || strings.TrimSpace(value) == "" {
		return defaultValue
	}

	return value
}

func AskYesOrNo(color string, message string, defaultValue bool, shouldPrompt bool) bool {
	if !shouldPrompt {
		return defaultValue
	}

	var prompt string
	if defaultValue {
		prompt = "[Y/n]"
	} else {
		prompt = "[y/N]"
	}

	answer := "none"
	for answer != "y" && answer != "n" && answer != "" {
		fmt.Print(color, fmt.Sprintf(" %s %s ", message, prompt), MsgReset)
		if _, err := fmt.Scanln(&answer); err != nil {
			return defaultValue
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
	}

	return (answer == "" && defaultValue) || answer == "y"
}

func AskMultipleChoices(message string, choices []string) string {
	answer := -1
	fmt.Println(MsgL1, fmt.Sprintf("%s", message), MsgReset)
	for answer < 1 || answer > len(choices) {
		for i, choice := range choices {
			fmt.Printf("    %d: %s\n", i+1, choice)
		}
		fmt.Printf(" > ")
		if _, err := fmt.Scanln(&answer); err != nil {
			fmt.Fprint(os.Stderr, "Not a valid integer\n")
		}
	}
	return choices[answer-1]
}
