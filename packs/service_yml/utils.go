package service_yml

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"gopkg.in/yaml.v2"
	"unicode"
	"github.com/cloud66/starter/common"
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

func generatePortsFromShortSyntax(shortSyntax string, clusterPorts []KubesPorts, nodePorts []KubesPorts, nodePort int) ([]KubesPorts, []KubesPorts, []KubesPorts, int) {
	var dPorts []KubesPorts
	containerStr := ""
	httpStr := ""
	httpsStr := ""
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
	if http == 0 && https == 0 {
		//crate new node of type ClusterIP
		clusterPorts = appendNewPortNoNodePort(clusterPorts, strconv.Itoa(container)+"-expose", container, container, "", 0)
		dPorts = appendNewPortNoNodePort(dPorts, strconv.Itoa(container)+"-expose", 0, 0, "", container)
	} else {
		if http != 0 {
			nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-http", container, http, "TCP", 0, nodePort)
			//create new port with name "http" of type NodePort
		}
		if https != 0 {
			//create new port with name "https of type NodePort
			nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-https", container, https, "TCP", 0, nodePort)
		}
		dPorts = appendNewPortNoNodePort(dPorts, strconv.Itoa(container)+"-tcp", 0, 0, "TCP", container)

	}
	return dPorts, clusterPorts, nodePorts, nodePort
}

func generatePortsFromLongSyntax(longSyntax ServicePort, clusterPorts []KubesPorts, nodePorts []KubesPorts, nodePort int) ([]KubesPorts, []KubesPorts, []KubesPorts, int) {
	var dPorts []KubesPorts

	container, http, https, tcp, udp := getIntFromServicePort(longSyntax)

	if tcp == 0 && http == 0 && https == 0 && udp == 0 {
		//type "ClusterIP"
		clusterPorts = appendNewPortNoNodePort(clusterPorts, strconv.Itoa(container)+"-expose", container, container, "", 0)
		dPorts = appendNewPortNoNodePort(dPorts, strconv.Itoa(container)+"-expose", 0, 0, "", container)
	} else if udp == 0 {
		if tcp != 0 {
			nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-tcp", container, tcp, "TCP", 0, nodePort)
		}
		if http != 0 {
			nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-http", container, http, "TCP", 0, nodePort)
		}
		if https != 0 {
			nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-https", container, https, "TCP", 0, nodePort)
		}
		dPorts = appendNewPortNoNodePort(dPorts, strconv.Itoa(container)+"-tcp", 0, 0, "TCP", container)

	} else if udp != 0 {
		nodePorts, nodePort = appendNewPort(nodePorts, strconv.Itoa(container)+"-udp", container, udp, "UDP", 0, nodePort)
		dPorts = appendNewPortNoNodePort(dPorts, strconv.Itoa(container)+"-udp", 0, 0, "UDP", container)

	}
	return dPorts, clusterPorts, nodePorts, nodePort
}

func generateServicesRequiredByPorts(serviceName string, serviceSpecs ServiceYMLService, nodePort int) ([]KubesPorts, []KubesService, int) {
	services := []KubesService{}
	var dPorts, cPorts, nPorts, clusterPorts, nodePorts, deployPorts []KubesPorts
	for _, v := range serviceSpecs.Ports {
		nPorts = []KubesPorts{}
		dPorts = []KubesPorts{}
		cPorts = []KubesPorts{}
		switch v.(type) {
		case map[interface{}]interface{}:
			var longSyntaxPort ServicePort
			temp, er := yaml.Marshal(v)
			CheckError(er)
			er = yaml.Unmarshal(temp, &longSyntaxPort)
			CheckError(er)
			deployPorts, clusterPorts, nodePorts, nodePort = generatePortsFromLongSyntax(longSyntaxPort, clusterPorts, nodePorts, nodePort)
		case string:
			shortSyntax := v.(string)
			deployPorts, clusterPorts, nodePorts, nodePort = generatePortsFromShortSyntax(shortSyntax, clusterPorts, nodePorts, nodePort)
		case int:
			shortSyntax := strconv.Itoa(v.(int))
			deployPorts, clusterPorts, nodePorts, nodePort = generatePortsFromShortSyntax(shortSyntax, clusterPorts, nodePorts, nodePort)
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
	if len(nodePorts) == 0 {
		if len(clusterPorts) > 0 {
			clusterService := generateService("ClusterIP", serviceSpecs, serviceName, clusterPorts)
			services = append(services, clusterService)
		}
	} else {
		nodeService := generateService("NodePort", serviceSpecs, serviceName, nodePorts)
		if len(clusterPorts) > 0 {
			clusterService := generateService("ClusterIP", serviceSpecs, serviceName, clusterPorts)
			services = append(services, clusterService)
		}
		services = append(services, nodeService)
	}

	return deployPorts, services, nodePort
}

func generateService(serviceType string, serviceSpecs ServiceYMLService, serviceName string, ports []KubesPorts) KubesService {
	service := KubesService{}

	service = KubesService{ApiVersion: "v1",
		Kind:                      "Service",
		Metadata: Metadata{
			Name:   serviceName + "-sv" + strings.ToLower(serviceType[:1]),
			Labels: serviceSpecs.Tags,
		},
		Spec: Spec{
			Type:  serviceType,
			Ports: ports,
		},
	}

	return service
}

func setDbDeploymentPorts(dbName string) []KubesPorts {
	ports, _ := getExposedPorts(dbName)
	return ports
}

func setDbServicePorts(dbName string) ([]KubesPorts) {
	_, ports := getExposedPorts(dbName)
	return ports
}

func getExposedPorts(dbName string) ([]KubesPorts, []KubesPorts) {
	dPorts := []KubesPorts{}
	sPorts := []KubesPorts{}

	switch dbName {
	case "mysql":
		dPorts = appendNewPortNoNodePort(dPorts, "mysql", 0, 0, "", 3306)
		sPorts = appendNewPortNoNodePort(sPorts, "mysql", 3306, 3306, "", 0)
	case "redis":
		dPorts = appendNewPortNoNodePort(dPorts, "redis", 0, 0, "", 6379)
		sPorts = appendNewPortNoNodePort(sPorts, "redis", 6379, 6379, "", 0)
	case "postgresql":
		dPorts = appendNewPortNoNodePort(dPorts, "postgresql", 0, 0, "", 5432)
		sPorts = appendNewPortNoNodePort(sPorts, "postgresql", 5432, 5432, "", 0)
	case "mongodb":
		dPorts = appendNewPortNoNodePort(dPorts, "mongodb", 0, 0, "", 27017)
		sPorts = appendNewPortNoNodePort(sPorts, "mongodb", 27017, 27017, "", 0)
	case "elasticsearch":
		dPorts = appendNewPortNoNodePort(dPorts, "elasticsearch", 0, 0, "", 9200)
		sPorts = appendNewPortNoNodePort(sPorts, "elasticsearch", 9200, 9200, "", 0)
		dPorts = appendNewPortNoNodePort(dPorts, "elasticsearch", 0, 0, "", 93000)
		sPorts = appendNewPortNoNodePort(sPorts, "elasticsearch", 9300, 9300, "", 0)
	case "glusterfs":
		dPorts = appendNewPortNoNodePort(dPorts, "glusterfs", 0, 0, "TCP", 24007)
		sPorts = appendNewPortNoNodePort(sPorts, "glusterfs", 24007, 24007, "TCP", 0)
		dPorts = appendNewPortNoNodePort(dPorts, "glusterfs", 0, 0, "UDP", 24008)
		sPorts = appendNewPortNoNodePort(sPorts, "glusterfs", 24008, 24008, "UDP", 0)
	case "influxdb":
		dPorts = appendNewPortNoNodePort(dPorts, "influxdb", 0, 0, "", 8086)
		sPorts = appendNewPortNoNodePort(sPorts, "influxdb", 8086, 8086, "", 0)
	case "rabbitmq":
		dPorts = appendNewPortNoNodePort(dPorts, "rabbitmq", 0, 0, "", 15672)
		sPorts = appendNewPortNoNodePort(sPorts, "rabbitmq", 15672, 15672, "", 0)
	default:
		common.PrintlnWarning("Not a recognized database.")
	}

	return dPorts, sPorts
}

func appendNewPort(ports []KubesPorts, name string, port int, targetPort int, protocol string, containerPort int, nodePort int) ([]KubesPorts, int) {
	ports = append(ports, KubesPorts{
		Name:          name + "-" + strconv.Itoa(nodePort),
		Port:          port,
		TargetPort:    targetPort,
		Protocol:      protocol,
		ContainerPort: containerPort,
		NodePort:      nodePort,
	})
	nodePort++
	return ports, nodePort
}

func appendNewPortNoNodePort(ports []KubesPorts, name string, port int, targetPort int, protocol string, containerPort int) ([]KubesPorts) {
	ports = append(ports, KubesPorts{
		Name:          name,
		Port:          port,
		TargetPort:    targetPort,
		Protocol:      protocol,
		ContainerPort: containerPort,
	})
	return ports
}

func getIntFromServicePort(longSyntax ServicePort) (int, int, int, int, int) {
	var container, http, https, tcp, udp int

	container = getIntFromVal(longSyntax.Container)
	http = getIntFromVal(longSyntax.Http)
	https = getIntFromVal(longSyntax.Https)
	tcp = getIntFromVal(longSyntax.Tcp)
	udp = getIntFromVal(longSyntax.Udp)

	return container, http, https, tcp, udp
}

func getIntFromVal(value string) int {
	var temp string
	var i int
	if len(value) > 0 {
		if value[0] == '"' {
			i = 1
		} else {
			i = 0
		}
		for ; i < len(value); i++ {
			if !unicode.IsDigit(rune(value[i])) {
				break
			} else {
				temp = string(append([]byte(temp), value[i]))
			}
		}
		if temp != "" {
			result, err := strconv.Atoi(temp)
			CheckError(err)
			return result
		}
	}
	return 0
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
