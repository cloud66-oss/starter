package common

import (
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"

	"github.com/blang/semver"
)

// Looks for node version in the package.json. If found returns true, version if not false, ""
func GetNodeVersion(packageJsonFile string) (bool, string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	if data["engines"] == nil {
		return false, ""
	}

	if nodeVersion, ok := data["engines"].(map[string]interface{})["node"].(string); ok {
		nodeVersion = strings.Split(nodeVersion, " ||")[0]
		nodeVersion = strings.Replace(strings.Trim(nodeVersion, "^>=~"), "x", "0", -1)
		v1, err := semver.Make(nodeVersion)
		if err != nil {
		  return true, nodeVersion
		}
		nodeVersion = fmt.Sprintf("%d.%d.%d", v1.Major, v1.Minor, v1.Patch)
		return true, nodeVersion
	}
	return false, ""
}

func GetNodeDatabase(packageJsonFile string, databaseNames ...string) (bool, string) {
	found, name := GetDependencyVersion(packageJsonFile, databaseNames...)
	return found, name
}

func GetDependencyVersion(packageJsonFile string, dependencyNames ...string) (bool, string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	for dependency, version := range data["dependencies"].(map[string]interface{}) {
		for _, dependencyName := range dependencyNames {
			found := dependencyName == dependency
			if found {
				return true, version.(string)
			}
		}

	}

	if data["optionalDependencies"] != nil {
		for dependency, version := range data["optionalDependencies"].(map[string]interface{}) {
			for _, dependencyName := range dependencyNames {
				found := dependencyName == dependency
				if found {
					return true, version.(string)
				}
			}

		}
	}



	return false, ""
}

func GetScriptsStart(packageJsonFile string) (bool, string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	if data["scripts"] == nil {
		return false, ""
	}

	if start, ok := data["scripts"].(map[string]interface{})["start"].(string); ok {
		return true, start
	} else {
		return false, ""
	}
}
