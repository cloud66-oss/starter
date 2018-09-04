package service_yml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ServiceYml struct {
	Services  map[string]Service `yaml:"services,omitempty"`
	Databases []string           `yaml:"databases,omitempty"`
}

func (s *ServiceYml) UnmarshalFromFile(path string) error {
	var err error
	_, err = os.Stat(path)
	CheckError(err)

	yamlFile, err := ioutil.ReadFile(path)

	serviceYml := ServiceYml{
		Services:  map[string]Service{},
		Databases: make([]string, 1),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &serviceYml); err != nil {
		fmt.Println(err.Error())
	}
	s.Services = serviceYml.Services
	s.Databases = serviceYml.Databases

	return nil
}

func (s ServiceYml) MarshalToFile(path string) error {
	file, err := yaml.Marshal(s)

	file = []byte("# Generated with <3 by Cloud66\n\n" + string(file))

	err = ioutil.WriteFile(path, file, 0644)
	CheckError(err)

	return nil
}
