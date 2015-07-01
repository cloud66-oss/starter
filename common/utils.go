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
	rubyVersionRegex = regexp.MustCompile("ruby\\s['\"](.*?)['\"]")

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

// Looks for ruby version in the gemfile. If found returns true, version if not false, ""
func GetRubyVersion(gemFile string) (bool, string) {
	buf, err := ioutil.ReadFile(gemFile)
	if err != nil {
		return false, err.Error()
	}

	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if rubyVersionRegex.MatchString(line) {
			sm := rubyVersionRegex.FindStringSubmatch(line)
			return true, sm[1]
		}
	}

	return false, ""
}

// returns bool = found any of the gems or not and string = the version of the first found
func GetGemVersion(gemFile string, gemNames ...string) (bool, string) {
	buf, err := ioutil.ReadFile(gemFile)
	if err != nil {
		return false, err.Error()
	}

	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		for _, gemName := range gemNames {
			found, version := ParseLineForGem(gemName, line)
			if found {
				return true, version
			}
		}
	}

	return false, ""

}

// Checks a line to see if it contains the given gem. returns true, version or false, ""
func ParseLineForGem(gemName string, line string) (bool, string) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		// empty or comment
		return false, ""
	}

	re := regexp.MustCompile(fmt.Sprintf("gem\\s['\"]%s['\"]\\s*,?\\s*(?P<version>['\"].*?['\"])?", gemName))
	if !re.MatchString(line) {
		return false, ""
	} else {
		sm := re.FindStringSubmatch(line)

		if len(sm) > 0 {

			result := strings.Replace(sm[1], "'", "", -1)
			result = strings.Replace(result, "\"", "", -1)

			return true, result
		} else {
			return true, ""
		}
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


