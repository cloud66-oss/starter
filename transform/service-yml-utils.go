package transform

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"gopkg.in/yaml.v2"
	"unicode"
	"github.com/cloud66/starter/common"
	"sort"
	"github.com/cloud66/starter/definitions/kubernetes"
	"github.com/cloud66/starter/definitions/service-yml"
)

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

func handleVolumes(serviceVolumes []string) []kubernetes.VolumeMounts {
	var kubeVolumes []kubernetes.VolumeMounts
	var outputWarning, readOnly bool

	for _, volume := range serviceVolumes {
		name := ""
		mountPath := ""
		var i int
		readOnly = false
		outputWarning = false

		if volume[0] == '"' {
			i = 1
		} else {
			i = 0
		}
		if volume[i] != '/' {
			outputWarning = true
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

		if outputWarning == true {
			common.PrintlnWarning("Path \"%s:%s\" not absolute! Please modify manually.", name, mountPath)
		}

		kubeVolume := kubernetes.VolumeMounts{
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

func generatePortsFromLongSyntax(longSyntax service_yml.Port, clusterPorts []kubernetes.Port, nodePorts []kubernetes.Port, nodePort int) ([]kubernetes.Port, []kubernetes.Port, []kubernetes.Port, int) {
	var dPorts []kubernetes.Port

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

func generateServicesRequiredByPorts(serviceName string, serviceSpecs service_yml.Service, nodePort int) ([]kubernetes.Port, []kubernetes.KubesService, int) {
	services := []kubernetes.KubesService{}
	var dPorts, cPorts, nPorts, clusterPorts, nodePorts, deployPorts []kubernetes.Port
	for _, v := range serviceSpecs.Ports {
		nPorts = []kubernetes.Port{}
		dPorts = []kubernetes.Port{}
		cPorts = []kubernetes.Port{}
		var longSyntaxPort service_yml.Port
		temp, er := yaml.Marshal(v)
		CheckError(er)
		er = yaml.Unmarshal(temp, &longSyntaxPort)
		CheckError(er)
		deployPorts, clusterPorts, nodePorts, nodePort = generatePortsFromLongSyntax(longSyntaxPort, clusterPorts, nodePorts, nodePort)

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

func generateService(serviceType string, serviceSpecs service_yml.Service, serviceName string, ports []kubernetes.Port) kubernetes.KubesService {
	service := kubernetes.KubesService{}

	service = kubernetes.KubesService{ApiVersion: "v1",
		Kind:                                 "Service",
		Metadata: kubernetes.Metadata{
			Name:   serviceName + "-sv" + strings.ToLower(serviceType[:1]),
			Labels: serviceSpecs.Tags,
		},
		Spec: kubernetes.Spec{
			Type:  serviceType,
			Ports: ports,
		},
	}

	return service
}

func setDbDeploymentPorts(dbName string) []kubernetes.Port {
	ports, _ := getExposedPorts(dbName)
	return ports
}

func setDbServicePorts(dbName string) ([]kubernetes.Port) {
	_, ports := getExposedPorts(dbName)
	return ports
}

func getExposedPorts(dbName string) ([]kubernetes.Port, []kubernetes.Port) {
	dPorts := []kubernetes.Port{}
	sPorts := []kubernetes.Port{}

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

func appendNewPort(ports []kubernetes.Port, name string, port int, targetPort int, protocol string, containerPort int, nodePort int) ([]kubernetes.Port, int) {
	ports = append(ports, kubernetes.Port{
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

func appendNewPortNoNodePort(ports []kubernetes.Port, name string, port int, targetPort int, protocol string, containerPort int) ([]kubernetes.Port) {
	ports = append(ports, kubernetes.Port{
		Name:          name,
		Port:          port,
		TargetPort:    targetPort,
		Protocol:      protocol,
		ContainerPort: containerPort,
	})
	return ports
}

func getIntFromServicePort(longSyntax service_yml.Port) (int, int, int, int, int) {
	var container, http, https, tcp, udp int

	container = longSyntax.Container
	http = longSyntax.Http
	https = longSyntax.Https
	tcp = longSyntax.Tcp
	udp = longSyntax.Udp

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

func composeWriter(file []byte, deployments []kubernetes.KubesDeployment, kubesServices []kubernetes.KubesService) []byte {
	var keys []string
	for _, k := range kubesServices {
		keys = append(keys, k.Metadata.Name)
	}
	sort.Strings(keys)
	indexPort := 31111

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

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
