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
		Http:      http,
		Https:     https,
	}
}

func stringToInt(strPort tempPort) Port {
	var port Port
	if strPort.Container != "" {
		if strPort.Container[0] == '"' {
			strPort.Container = strPort.Container[0:]
			strPort.Container = strPort.Container[:len(strPort.Container)-1]
		}
		port.Container, _ = strconv.Atoi(strPort.Container)

	}
	if strPort.Tcp != "" {
		if strPort.Tcp[0] == '"' {
			strPort.Tcp = strPort.Tcp[0:]
			strPort.Tcp = strPort.Tcp[:len(strPort.Tcp)-1]
		}
		port.Tcp, _ = strconv.Atoi(strPort.Tcp)

	}
	if strPort.Http != "" {
		if strPort.Http[0] == '"' {
			strPort.Http = strPort.Http[0:]
			strPort.Http = strPort.Http[:len(strPort.Http)-1]
		}
		port.Http, _ = strconv.Atoi(strPort.Http)

	}
	if strPort.Https != "" {
		if strPort.Https[0] == '"' {
			strPort.Https = strPort.Https[0:]
			strPort.Https = strPort.Https[:len(strPort.Https)-1]
		}
		port.Https, _ = strconv.Atoi(strPort.Https)

	}
	if strPort.Udp != "" {
		if strPort.Udp[0] == '"' {
			strPort.Udp = strPort.Udp[0:]
			strPort.Udp = strPort.Udp[:len(strPort.Udp)-1]
		}
		port.Udp, _ = strconv.Atoi(strPort.Udp)
	}
	return port
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
