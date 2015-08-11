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

	MsgReset string = ansi.ColorCode("reset")

	PrintTitle, PrintlnTitle     = printers(" ", ansi.ColorCode("green+h"))
	PrintL0, PrintlnL0           = printers(" ", ansi.ColorCode("magenta"))
	PrintL1, PrintlnL1           = printers(" ", ansi.ColorCode("white"))
	PrintL2, PrintlnL2           = printers(" ----> ", ansi.ColorCode("black+h"))
	PrintError, PrintlnError     = printers(" ", ansi.ColorCode("red"))
	PrintWarning, PrintlnWarning = printers(" ", ansi.ColorCode("yellow"))
)

func printers(prefix string, color string) (print func(format string, a ...interface{}), println func(format string, a ...interface{})) {
	print = func(format string, a ...interface{}) {
		fmt.Printf(color+prefix+format+MsgReset, a...)
	}
	println = func(format string, a ...interface{}) {
		print(format+"\n", a...)
	}
	return print, println
}

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
		PrintL1("%s: ", message)
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

	PrintL1("%s [%s] ", message, printedDefaultValue)
	value := ""
	if _, err := fmt.Scanln(&value); err != nil || strings.TrimSpace(value) == "" {
		return defaultValue
	}

	return value
}

func AskYesOrNo(message string, defaultValue bool, shouldPrompt bool) bool {
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
		PrintL1("%s %s ", message, prompt)
		if _, err := fmt.Scanln(&answer); err != nil {
			return defaultValue
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
	}

	return (answer == "" && defaultValue) || answer == "y"
}

func AskMultipleChoices(message string, choices []string) string {
	answer := -1
	PrintL1(message)
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
