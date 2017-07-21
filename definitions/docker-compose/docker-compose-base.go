package docker_compose

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type DockerCompose struct {
	Services map[string]Service
	Version  string
}

func (d DockerCompose) UnmarshalFromFile(path string) error {
	var err error
	_, err = os.Stat(path)
	CheckError(err)

	yamlFile, err := ioutil.ReadFile(path)

	dockerCompose := DockerCompose{
		Services: make(map[string]Service),
		Version:  "",
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &dockerCompose); err != nil {
		fmt.Println(err.Error())
	}

	if len(dockerCompose.Services) == 0 {
		err = yaml.Unmarshal([]byte(yamlFile), &dockerCompose.Services)
		CheckError(err)
	}

	d.Services = dockerCompose.Services
	d.Version = dockerCompose.Version

	return nil
}

func (d DockerCompose) MarshalToFile(path string) error {
	file, err := yaml.Marshal(d)
	file = []byte("# Generated with <3 by Cloud66\n\n" + string(file))

	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
