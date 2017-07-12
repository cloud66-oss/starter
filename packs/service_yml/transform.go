package service_yml

import (
	"github.com/cloud66/starter/common"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

func Transformer(filename string, formatTarget string) error {

	var err error
	_, err = os.Stat(formatTarget)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(formatTarget)
		CheckError(err)
		defer file.Close()
	} else {
		common.PrintError("File %s already exists. Will be overwritten.\n", formatTarget)
	}

	yamlFile, err := ioutil.ReadFile(filename)

	serviceYML := ServiceYml{
		Services: make(map[string]ServiceYMLService),
		Dbs:      []string{},
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &serviceYML); err != nil {
		fmt.Println(err.Error())
	}

	file := copyToKubes(serviceYML)

	err = ioutil.WriteFile(formatTarget, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func copyToKubes(serviceYml ServiceYml) []byte {

	var file []byte
	var deploy KubesDeployment

	file = []byte("# Generated with <3 by Cloud66\n\n")

	for _, dbName := range serviceYml.Dbs {
		file = []byte(string(file) + "####### " + string(dbName) + " #######" + "\n\n")

		service := KubesService{
			ApiVersion: "extenstions/v1beta1",
			Kind:       "Service",
			Metadata: Metadata{
				Name:   dbName + "-svc",
				Labels: map[string]string{},
			},
			Spec: Spec{
				Type:  "ClusterIP",
				Ports: setDbServicePorts(dbName),
			},
		}
		deploy := KubesService{ApiVersion: "extensions/v1beta1",
			Kind:                      "Deployment",
			Metadata: Metadata{
				Name: dbName + "-deployment",
			},
			Spec: Spec{
				Template: Template{
					Metadata: Metadata{
						Labels: service.Metadata.Labels,
					},
					PodSpec: PodSpec{
						Containers: []Containers{
							{
								Name:  dbName,
								Image: dbName + ":latest",
								Ports: setDbDeploymentPorts(dbName),
							},
						},
					},
				},
			},
		}

		//write db service
		fileServices, er := yaml.Marshal(service)
		CheckError(er)
		file = []byte(string(file) + string(handleEnvVarsFormat(fileServices)) + "---\n")

		//write db deployment
		fileDeployments, er := yaml.Marshal(deploy)
		CheckError(er)
		file = []byte(string(file) + string(handleEnvVarsFormat(fileDeployments)) + "---\n")

	}

	for serviceName, serviceSpecs := range serviceYml.Services {

		//gets ports to populate deployment and generates the required service(s)
		deployPorts, services := handlePorts(serviceName, serviceSpecs)

		deploy = KubesDeployment{ApiVersion: "extensions/v1beta1",
			Kind:                        "Deployment",
			Metadata: Metadata{
				Name: serviceName + "-deployment",
			},
			Spec: Spec{
				Template: Template{
					Metadata: Metadata{
						Labels: serviceSpecs.Tags,
					},
					PodSpec: PodSpec{
						Containers: []Containers{
							{
								Name:       serviceName,
								Image:      serviceSpecs.Image,
								Command:    serviceSpecs.Command,
								Ports:      deployPorts,
								WorkingDir: serviceSpecs.WorkDir,
								SecurityContext: SecurityContext{
									Priviliged: serviceSpecs.Privileged,
								},
								Lifecycle: Lifecycle{
									PostStart: Handler{
										Exec: Exec{
											Command: serviceSpecs.PostStartCommand,
										},
									},
									PreStop: Handler{
										Exec: Exec{
											Command: serviceSpecs.PreStopCommand,
										},
									},
								},
								Resources: KubesResources{
									Limits: Limits{
										Cpu:    serviceSpecs.Constraints.Resources.Cpu,
										Memory: serviceSpecs.Constraints.Resources.Memory,
									},
								},
							},
						},
					},
				},
			},
		}

		kubeVolumes := handleVolumes(serviceSpecs.Volumes)
		deploy.Spec.Template.PodSpec.Containers[0].VolumeMounts = kubeVolumes

		keys, values := getKeysValues(serviceSpecs.EnvVars)
		if len(keys) > 0 {
			for k := 0; k < len(keys); k++ {
				if values[k] == "\"\"" {
					values[k] = ""
				}
				env := EnvVar{
					Name:  keys[k],
					Value: values[k],
				}
				deploy.Spec.Template.PodSpec.Containers[0].Env = append(deploy.Spec.Template.PodSpec.Containers[0].Env, env)
			}
		}
		file = []byte(string(file) + "####### " + string(serviceName) + " #######" + "\n\n")
		for _, service := range services {
			fileServices, er := yaml.Marshal(service)
			CheckError(er)
			file = []byte(string(file) + string(handleEnvVarsFormat(fileServices)) + "---\n")
		}

		fileDeployments, er := yaml.Marshal(deploy)
		CheckError(er)

		file = []byte(string(file) + string(handleEnvVarsFormat(fileDeployments)) + "---\n")

		//delete the last row of "---"
		if len(file) > 3 {
			file = file[:len(file)-4]
		}
	}

	return file
}
