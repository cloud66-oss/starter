package docker_compose

import (
	"unicode"
	"fmt"
	"os"
	"strconv"
	"strings"
	"bufio"
)

func getKeyValue(line string) (string, string) {
	var key, value string
	var k int
	for k = 0; k < len(line); k++ {
		if !unicode.IsSpace(rune(line[k])) {
			break
		}
	}
	for ; k < len(line); k++ {
		if line[k] == '=' {
			break
		} else {
			key = string(append([]byte(key), line[k]))
		}
	}
	for k = k + 1; k < len(line); k++ {
		if line[k] == '\n' {
			break
		} else {
			value = string(append([]byte(value), line[k]))
		}
	}

	return key, value
}

func shortPortToLong(str string) Port {
	var port Port
	var i int
	var host, container, protocol string
	if str[0] == '"' {
		i = 1
	} else {
		i = 0
	}
	for ; i < len(str); i++ {
		if unicode.IsDigit(rune(str[i])) {
			break
		} else {
			host = string(append([]byte(host), str[i]))
		}
	}
	for i = i + 1; i < len(str); i++ {
		if !unicode.IsDigit(rune(str[i])) {
			break
		} else {
			container = string(append([]byte(container), str[i]))
		}
	}

	protocol = "tcp"
	if i<len(str)-1{
		if strings.Contains(str[i:], "udp"){
			protocol="udp"
		}
	}

	target, err := strconv.Atoi(container)
	CheckError(err)
	published, err := strconv.Atoi(host)
	CheckError(err)

	port = Port{
		Protocol:  protocol,
		Target:    target,
		Published: published,
	}

	return port
}

func shortSecretToLong(str string) Secret {
	var secret Secret
	secret.Source = str
	return secret
}

func shortVolumeToLong(str string) Volume {
	var volume Volume
	var i int
	var source, target string
	var readOnly bool

	if str[0] == '"' {
		i = 1
	} else {
		i = 0
	}
	for ; i < len(str); i++ {
		if str[i] == ':' {
			break
		} else {
			source = string(append([]byte(source), str[i]))
		}
	}
	for i = i + 1; i < len(str); i++ {
		if str[i] == ':' || str[i] == '\n' || str[i] == '"' {
			break
		} else {
			target = string(append([]byte(target), str[i]))
		}
	}
	if i < len(str)-2 {
		if strings.Contains(str[i+1:], "ro") {
			readOnly = true
		}
	}
	volume = Volume{
		Type:     "volume",
		Source:   source,
		Target:   target,
		ReadOnly: readOnly,
	}
	return volume
}


func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
