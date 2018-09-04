package common

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func GetPythonVersion() (hasFound bool, version string) {
	cmd := exec.Command("python", "--version")
	b, err := cmd.CombinedOutput() // NOTE: cmd.Output() would not work here because the Python version is written on stderr
	if err != nil {
		return false, ""
	}
	return true, strings.TrimSpace(strings.Split(string(b), " ")[1])
}

func IsDjangoProject(rootDir string) bool {
	return FileExists(filepath.Join(rootDir, "manage.py"))
}

func PythonPackages(requirementsTxt string) ([]string, error) {
	packageRegexp := regexp.MustCompile(`(?m)^(\w+)`)
	includeRegexp := regexp.MustCompile(`(?m)^-r[[:blank:]]+(.*)[[:blank:]]*$`)

	content, err := ioutil.ReadFile(requirementsTxt)
	if err != nil {
		return nil, err
	}

	var packages []string
	for _, match := range packageRegexp.FindAllStringSubmatch(string(content), -1) {
		packages = append(packages, match[1])
	}

	for _, match := range includeRegexp.FindAllStringSubmatch(string(content), -1) {
		otherPackages, err := PythonPackages(match[1])
		if err != nil {
			return nil, err
		}
		packages = append(packages, otherPackages...)
	}

	return packages, nil
}
