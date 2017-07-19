package service_yml

import (
	"strconv"
	"fmt"
	"os"
)

func shortPortToLong(shortSyntax string) Port {
	var containerStr, httpStr, httpsStr string
	var i int

	if shortSyntax[0] == '"' {
		i = 1
	} else {
		i = 0
	}

	for ; i < len(shortSyntax); i++ {
		if shortSyntax[i] == ':' {
			break
		} else {
			containerStr = string(append([]byte(containerStr), shortSyntax[i]))
		}
	}

	for i = i + 1; i < len(shortSyntax); i++ {
		if shortSyntax[i] == ':' {
			break
		} else {
			httpStr = string(append([]byte(httpStr), shortSyntax[i]))
		}
	}

	for i = i + 1; i < len(shortSyntax); i++ {
		if shortSyntax[i] == '"' || shortSyntax[i] == '\n' {
			break
		} else {
			httpsStr = string(append([]byte(httpsStr), shortSyntax[i]))
		}
	}

	var container, http, https int
	var err error
	if containerStr != "" {
		container, err = strconv.Atoi(containerStr)
		CheckError(err)
	} else {
		container = 0
	}
	if httpStr != "" {
		http, err = strconv.Atoi(httpStr)
		CheckError(err)
	} else {
		http = 0
	}
	if httpsStr != "" {
		https, err = strconv.Atoi(httpsStr)
		CheckError(err)
	} else {
		https = 0
	}

	return Port{
		Container: container,
		Http: http,
		Https: https,
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
