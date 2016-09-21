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
		  return false, ""
		}
		nodeVersion := GetClosedAllowedNodeVersion(v1.Major, v1.Minor, v1.Patch)
		return true, nodeVersion
	}
	return false, ""
}

func GetClosedAllowedNodeVersion(major uint64, minor uint64, patch uint64) (string) {
	for _, version := range allowedNodeVersions {
		if strings.Index(version, fmt.Sprintf("%d.%d", major, minor)) == 0 {
			return version
		}
	}
	//last resort
	return "latest"
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

func SetAllowedNodeVersions(versions []string) {
	allowedNodeVersions = versions
}

func GetAllowedNodeVersions() []string {
	return allowedNodeVersions
}

var allowedNodeVersions = []string { 
		"0.10.28",
        "0.10.30",
        "0.10.31",
        "0.10.32",
        "0.10.33",
        "0.10.34",
        "0.10.35",
        "0.10.36",
        "0.10.37",
        "0.10.38",
        "0.10.39",
        "0.10.40",
        "0.10.41",
        "0.10.42",
        "0.10.43",
        "0.10.44",
        "0.10.45",
        "0.10.46",
        "0.10",
        "0.11.13",
        "0.11.14",
        "0.11.15",
        "0.11.16",
        "0.11",
        "0.12.0",
        "0.12.1",
        "0.12.10",
        "0.12.11",
        "0.12.12",
        "0.12.13",
        "0.12.14",
        "0.12.15",
        "0.12.2",
        "0.12.3",
        "0.12.4",
        "0.12.5",
        "0.12.6",
        "0.12.7",
        "0.12.8",
        "0.12.9",
        "0.12",
        "0.8.28",
        "0.8",
        "0",
        "4.0.0",
        "4.0",
        "4.1.0",
        "4.1.1",
        "4.1.2",
        "4.1",
        "4.2.0",
        "4.2.1",
        "4.2.2",
        "4.2.3",
        "4.2.4",
        "4.2.5",
        "4.2.6",
        "4.2",
        "4.3.0",
        "4.3.1",
        "4.3.2",
        "4.3",
        "4.4.0",
        "4.4.1",
        "4.4.2",
        "4.4.3",
        "4.4.4",
        "4.4.5",
        "4.4.6",
        "4.4.7",
        "4.4",
        "4.5.0",
        "4.5",
        "4",
        "5.0.0",
        "5.0",
        "5.1.0",
        "5.1.1",
        "5.1",
        "5.10.0",
        "5.10.1",
        "5.10",
        "5.11.0",
        "5.11.1",
        "5.11",
        "5.12.0",
        "5.12",
        "5.2.0",
        "5.2",
        "5.3.0",
        "5.3",
        "5.4.0",
        "5.4.1",
        "5.4",
        "5.5.0",
        "5.5",
        "5.6.0",
        "5.6",
        "5.7.0",
        "5.7.1",
        "5.7",
        "5.8.0",
        "5.8",
        "5.9.0",
        "5.9.1",
        "5.9",
        "5",
        "6.0.0",
        "6.0",
        "6.1.0",
        "6.1",
        "6.2.0",
        "6.2.1",
        "6.2.2",
        "6.2",
        "6.3.0",
        "6.3.1",
        "6.3",
        "6.4.0",
        "6.4",
        "6.5.0",
        "6.5",
        "6.6.0",
        "6.6",
        "6",
    }