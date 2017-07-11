package service_yml

import (
	"github.com/cloud66/starter/common"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"strconv"
)

func Transformer(filename string, formatTarget string, gitURL string, gitBranch string, shouldPrompt bool) error {

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

	kubes := Kubes{
		Services:    []KubesService{},
		Deployments: []KubesDeployment{},
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &serviceYML); err != nil {
		fmt.Println(err.Error())
	}

	kubes = copyToKubes(serviceYML, shouldPrompt)

	file, err := yaml.Marshal(kubes)

	//Might need some formating for the "file"
	file = []byte("# Generated with <3 by Cloud66\n\n" + handleEnvVars(file))

	err = ioutil.WriteFile(formatTarget, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func copyToKubes(serviceYml ServiceYml, shouldPrompt bool) Kubes {

	kubes := Kubes{
		Services:    []KubesService{},
		Deployments: []KubesDeployment{},
	}

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
								//add some volumes here
								WorkingDir: serviceSpecs.WorkDir,
								Resources: KubesResources{
									Limits: Limits{
										Cpu:    strconv.Itoa(serviceSpecs.Constraints.Resources.Cpu),
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

		kubes.Services = append(kubes.Services, service)
		kubes.Deployments = append(kubes.Deployments, deploy)
	}

	return kubes
}
