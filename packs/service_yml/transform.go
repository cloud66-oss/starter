package service_yml

import (
	"github.com/cloud66/starter/common"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

func Transformer(filename string, formatTarget string, shouldPrompt bool) error {

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

	file := copyToKubes(serviceYML, shouldPrompt)

	err = ioutil.WriteFile(formatTarget, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func copyToKubes(serviceYml ServiceYml, shouldPrompt bool) []byte {

	kubes := Kubes{
		Services:    []KubesService{},
		Deployments: []KubesDeployment{},
	}
	var file []byte

	file = []byte("# Generated with <3 by Cloud66\n\n")

	for serviceName, serviceSpecs := range serviceYml.Services {

		deploy := KubesDeployment{ApiVersion: "extensions/v1beta1",
			Kind:                         "Deployment",
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
								Name:    serviceName,
								Image:   serviceSpecs.Image,
								Command: serviceSpecs.Command,
								//add some ports here
								WorkingDir: serviceSpecs.WorkDir,
								Resources: KubesResources{
									Limits: Limits{
										Cpu:    serviceSpecs.Constraints.Resources.Cpu,
										Memory: serviceSpecs.Constraints.Resources.Memory,
									},
								},
								SecurityContext: SecurityContext{
									Priviliged: serviceSpecs.Privileged,
								},
							},
						},
					},
				},
			},
		}
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
		kubeVolumes := handleVolumes(serviceSpecs.Volumes)
		deploy.Spec.Template.PodSpec.Containers[0].VolumeMounts = kubeVolumes

		service := KubesService{ApiVersion: "extensions/v1beta1",
			Kind:                       "Service",
			Metadata: Metadata{
				Name:   serviceName + "-svc",
				Labels: serviceSpecs.Tags,
			},
			Spec: Spec{
				//add some ports here
			},
		}

		fileServices, er := yaml.Marshal(service)
		CheckError(er)
		fileDeployments, er := yaml.Marshal(deploy)
		CheckError(er)

		file = []byte(string(file) + string(handleEnvVarsFormat(fileServices)) + "---\n" + string(handleEnvVarsFormat(fileDeployments)) + "---\n")
		kubes.Services = append(kubes.Services, service)
		kubes.Deployments = append(kubes.Deployments, deploy)
	}
	
	//delete the last row of "---"
	if len(file) > 3 {
		file = file[:len(file)-3]
	}
	return file
}
