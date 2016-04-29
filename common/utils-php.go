package common

import (
	// "fmt"
	"io/ioutil"
	"encoding/json"
	//"strings"
	"regexp"
)

// Looks for node version in the package.json. If found returns true, version if not false, ""
func GetPHPVersion(composerJSONfile string) (bool, string) {
	buf, err := ioutil.ReadFile(composerJSONfile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	if data["require"] == nil {
		return false, ""
	}

	if phpVersion, ok := data["require"].(map[string]interface{})["php"].(string); ok {
		re := regexp.MustCompile("[0-9][.][0-9]")
		return true, re.FindString(phpVersion)
	} else {
		return false, ""
	}

}

func GetPHPDatabase(composerJSONfile string, databaseName string) (bool, string) {
	//found, name := GetDependencyVersion(composerJSONfile, databaseName)
	return false, ""
}

