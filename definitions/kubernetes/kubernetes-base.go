package kubernetes

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/cloud66-oss/starter/definitions/docker-compose"
)

type Kubernetes struct {
	Services    []KubesService
	Deployments []KubesDeployment
}

func (k Kubernetes) UnmarshalFromFile(path string) error {
	//needs changes if used
	var err error
	_, err = os.Stat(path)
	docker_compose.CheckError(err)

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

	file := []byte("# Generated with <3 by Cloud66\n\n")
	file = composeWriter(file, k.Deployments, k.Services)

	err := ioutil.WriteFile(path, file, 0644)
	docker_compose.CheckError(err)

	return nil
}
