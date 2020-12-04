package kubernetes

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func composeWriter(file []byte, deployments []KubesDeployment, kubesServices []KubesService) []byte {
	var keys []string
	for _, k := range kubesServices {
		keys = append(keys, k.Metadata.Name)
	}
	sort.Strings(keys)
	indexPort := 31111 //magic number to hope it will not overlap with other unrelated ports

	for _, k := range keys {
		for i := 0; i < len(kubesServices); i++ {
			if kubesServices[i].Metadata.Name == k {
				if len(kubesServices[i].Spec.Ports) > 0 {
					for v, port := range kubesServices[i].Spec.Ports {
						if port.NodePort != 0 {
							port.NodePort = indexPort
							port.Name = port.Name[:len(port.Name)-5] + strconv.Itoa(indexPort)
							indexPort++
							kubesServices[i].Spec.Ports[v] = port
						}
					}
				}
				fileServices, err := yaml.Marshal(kubesServices[i])
				CheckError(err)
				file = []byte(string(file) + "####### " + strings.ToUpper(string(kubesServices[i].Metadata.Name)) + " - Service #######\n" + "\n" + string(finalFormat(fileServices)) + "---\n")
				break
			}
		}
	}

	keys = []string{}
	for _, k := range deployments {
		keys = append(keys, k.Metadata.Name)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for i := 0; i < len(deployments); i++ {
			if deployments[i].Metadata.Name == k {
				fileDeployments, err := yaml.Marshal(deployments[i])
				CheckError(err)
				file = []byte(string(file) + "---\n####### " + strings.ToUpper(string(deployments[i].Metadata.Name)) + " #######\n" + string(finalFormat(fileDeployments)))
				break
			}
		}
	}
	return file
}

func finalFormat(file []byte) string {
	finalFormat := ""

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "_env") {
			for j := 0; j < len(lines[i])-4; j++ {
				if lines[i][j] == '_' && lines[i][j+1] == 'e' && lines[i][j+2] == 'n' && lines[i][j+3] == 'v' {
					lines[i] = lines[i][:j] + "$" + lines[i][j+4:]
				}
			}
		}

		// handle empty value for env_vars
		if strings.Contains(lines[i], "value: ") && strings.Contains(lines[i], "\"\"") {
			for j := 0; j < len(lines[i]); j++ {
				if lines[i][j] == '"' {
					lines[i] = lines[i][:j]
					break
				}
			}
		}

		//format comment for empty image field
		if strings.Contains(lines[i], "image: '#") {
			for j := 0; j < len(lines[i]); j++ {
				if string(lines[i][j]) == "'" && j < len(lines[i])-1 {
					lines[i] = lines[i][:j] + " " + lines[i][j+1:]
				} else if string(lines[i][j]) == "'" {
					lines[i] = lines[i][:j]
				}
			}
		}

		finalFormat = finalFormat + lines[i] + "\n"
	}

	return finalFormat
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
