package service_yml

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"github.com/cloud66/starter/common"
	"gopkg.in/yaml.v2"
)

func handleEnvVarsFormat(file []byte) string {
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
		finalFormat = finalFormat + lines[i] + "\n"
	}

	return finalFormat
}

func handleVolumes(serviceVolumes []string) []VolumeMounts {
	var kubeVolumes []VolumeMounts

	for _, volume := range serviceVolumes {
		name := ""
		mountPath := ""
		var i int
		var readOnly bool
		if volume[0] == '"' {
			i = 1
		} else {
			i = 0
		}
		for ; i < len(volume); i++ {
			if volume[i] == ':' {
				break
			} else {
				name = string(append([]byte(name), volume[i]))
			}
		}

		for i = i + 1; i < len(volume); i++ {
			if volume[i] == ':' || volume[i] == '"' || volume[i] == '\n' {
				break
			} else {
				mountPath = string(append([]byte(mountPath), volume[i]))
			}
		}
		if i < len(volume)-2 {
			if volume[i] == ':' && volume[i+1] == 'r' && volume[i+2] == 'o' {
				readOnly = true
			}
		}
		kubeVolume := VolumeMounts{
			Name:      name,
			MountPath: mountPath,
			ReadOnly:  readOnly,
		}
		kubeVolumes = append(kubeVolumes, kubeVolume)
	}

	return kubeVolumes
}

func getKeysValues(env_vars map[string]string) ([]interface{}, []interface{}) {
	keys := []interface{}{}
	values := []interface{}{}
	for k, v := range env_vars {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func generatePortsFromShortSyntax(shortSyntax string, clusterPorts []KubesPorts, nodePorts []KubesPorts) ([]KubesPorts, []KubesPorts, []KubesPorts) {
	var dPort []KubesPorts
	x := ""
	y := ""
	z := ""
	var i int

	for i = 0; i < len(shortSyntax); i++ {
		if shortSyntax[i] == ':' {
			break
		} else {
			x = string(append([]byte(x), shortSyntax[i]))
		}
	}

	for i = i + 1; i < len(shortSyntax); i++ {
		if shortSyntax[i] == ':' {
			break
		} else {
			y = string(append([]byte(y), shortSyntax[i]))
		}
	}

	for i = i + 1; i < len(shortSyntax); i++ {
		if shortSyntax[i] == '"' || shortSyntax[i] == '\n' {
			break
		} else {
			z = string(append([]byte(z), shortSyntax[i]))
		}
	}

	if y == "" && z == "" {
		//crate new node of type ClusterIP
	}
	if y != "" {
		//create new port with name "http" of type NodePort
	}
	if z != "" {
		//create new port with name "https of type NodePort
	}

	return dPort, clusterPorts, nodePorts
}

func generatePortsFromLongSyntax(longSyntax ServicePort, clusterPorts []KubesPorts, nodePorts []KubesPorts) ([]KubesPorts, []KubesPorts, []KubesPorts) {
	var dPorts []KubesPorts

	if longSyntax.Tcp == "" && longSyntax.Https == "" && longSyntax.Https == "" && longSyntax.Udp == "" {
		//type "ClusterIP"
		cPort := KubesPorts{
			Port:       longSyntax.Container,
			TargetPort: longSyntax.Container,
		}
		dPort := KubesPorts{
			ContainerPort: longSyntax.Container,
		}

		clusterPorts = append(clusterPorts, cPort)
		dPorts = append(dPorts, dPort)
	}
	if longSyntax.Udp != "" {
		
	}
	if longSyntax.Tcp != "" {

	}
	if longSyntax.Http != "" {

	}
	if longSyntax.Https != "" {

	}

	return dPorts, clusterPorts, nodePorts
}

func handlePorts(serviceName string, serviceSpecs ServiceYMLService) ([]KubesPorts, []KubesService) {
	deployPorts := []KubesPorts{}
	services := []KubesService{}
	var dPorts, cPorts, nPorts, clusterPorts, nodePorts []KubesPorts
	for _, v := range serviceSpecs.Ports {
		common.PrintlnTitle("%#v", v)
		switch vv := v.(type) {
		case string:
			common.PrintlnL0("It s a string %v", vv)
			shortSyntax := v.(string)
			dPorts, cPorts, nPorts = generatePortsFromShortSyntax(shortSyntax, clusterPorts, nodePorts)
		case int:
			common.PrintlnL0("It s int %v", vv)
			shortSyntax := strconv.Itoa(v.(int))
			dPorts, cPorts, nPorts = generatePortsFromShortSyntax(shortSyntax, clusterPorts, nodePorts)
		case map[interface{}]interface{}:
			var longSyntaxPort ServicePort
			temp, er := yaml.Marshal(v)
			CheckError(er)
			er = yaml.Unmarshal(temp, &longSyntaxPort)
			CheckError(er)
			common.PrintlnTitle("It s a long one %v", vv)
			dPorts, cPorts, nPorts = generatePortsFromLongSyntax(longSyntaxPort, clusterPorts, nodePorts)
		}

		for _, port := range dPorts {
			deployPorts = append(deployPorts, port)
		}
		for _, port := range cPorts {
			clusterPorts = append(clusterPorts, port)
		}
		for _, port := range nPorts {
			nodePorts = append(nodePorts, port)
		}
	}

	//generate services with the specific type required by the found nodes
	clusterService := generateService("ClusterIP", serviceSpecs, serviceName, clusterPorts)
	nodeService := generateService("NodePorts", serviceSpecs, serviceName, nodePorts)
	services = append(services, clusterService)
	services = append(services, nodeService)

	return deployPorts, services
}

func generateService(serviceType string, serviceSpecs ServiceYMLService, serviceName string, ports []KubesPorts) KubesService {
	service := KubesService{}
	service = KubesService{ApiVersion: "extensions/v1beta1",
		Kind:                      "Service",
		Metadata: Metadata{
			Name:   serviceName + "-svc",
			Labels: serviceSpecs.Tags,
		},
		Spec: Spec{
			Type:  serviceType,
			Ports: ports,
		},
	}

	return service
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
