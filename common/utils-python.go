package common

import (
	"os/exec"
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
