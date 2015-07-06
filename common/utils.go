package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/mgutz/ansi"
)

var (
	procfileRegex    = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")
	envVarRegex      = regexp.MustCompile("\\$([A-Z_]+[A-Z0-9_]*)")

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

func LocalGitBranch(folder string) string {
	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", folder), "name-rev", "--name-only", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func RemoteGitUrl(folder string) string {
	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", folder), "config", "remote.origin.url").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func AskUser(message string, default_value string) string {
	fmt.Print(MsgL1, fmt.Sprintf(" %s [%s] ", message, default_value), MsgReset)
	value := ""
	if _, err := fmt.Scanln(&value); err != nil || strings.TrimSpace(value) == "" {
		return default_value
	}

	return value
}


