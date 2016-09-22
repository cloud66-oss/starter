package common

import (
	"fmt"
	"strings"
	"strconv"
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
			//no semver 
			i, err := strconv.Atoi(nodeVersion)
			if err != nil {
				return false, ""	
			}
			nodeVersion = GetClosedAllowedNodeVersion(uint64(i), 0, 0)
			return true, nodeVersion

		}
		nodeVersion = GetClosedAllowedNodeVersion(v1.Major, v1.Minor, v1.Patch)
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
	for _, version := range allowedNodeVersions {
		if strings.Index(version, fmt.Sprintf("%d", major)) == 0 {
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

	if _, ok := data["scripts"].(map[string]interface{})["start"].(string); ok {
		return true, "npm start"
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
	    "0.10.46",
        "0.12.15",
        "4.5.0",
        "5.12.0",
        "6.6.0",
    }
