package kubernetes

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/cloud66/starter/packs/service_yml"
)

type Kubernetes struct {
	Services    []KubesService
	Deployments []KubesDeployment
}

func (k Kubernetes) UnmarshalFromFile(path string) error {
	//needs changes if used
	var err error
	_, err = os.Stat(path)
	service_yml.CheckError(err)

	yamlFile, err := ioutil.ReadFile(path)

	kubernetes := Kubernetes{
		Services:    []KubesService{},
		Deployments: []KubesDeployment{},
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &kubernetes); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (k Kubernetes) MarshalToFile(path string) error {
	fileServices, err := yaml.Marshal(k.Services)
	service_yml.CheckError(err)
	fileDeployment, err := yaml.Marshal(k.Deployments)

	file := []byte("# Generated with <3 by Cloud66\n\n" + string(fileServices) + "#####\n   Deployments\n#####\n\n" + string(fileDeployment))

	err = ioutil.WriteFile(path, file, 0644)
	service_yml.CheckError(err)

	return nil
}
