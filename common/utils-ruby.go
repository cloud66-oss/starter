package common

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	rubyVersionRegex = regexp.MustCompile("ruby\\s['\"](.*?)['\"]")
)

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
