package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// returns bool = found any of the gems or not and string = the version of the first found
func GetGemVersion(gemFile string, gemNames ...string) (bool, string) {
	buf, err := ioutil.ReadFile(gemFile)
	if err != nil {
		return false, err.Error()
	}

	for _, gemName := range gemNames {
		re := regexp.MustCompile(fmt.Sprintf("[^#]gem\\s['\"]%s['\"]\\s*,?\\s*(?P<version>['\"].*?['\"])?", gemName))

		if !re.Match(buf) {
			return false, ""
		} else {
			sm := re.FindStringSubmatch(string(buf))

			if len(sm) > 0 {

				result := strings.Replace(sm[1], "'", "", -1)
				result = strings.Replace(result, "\"", "", -1)

				return true, result
			} else {
				return true, ""
			}
		}
	}

	return false, ""
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
