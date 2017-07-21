package service_yml

import (
	"os"
	"github.com/cloud66/starter/packs/service_yml"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type ServiceYml struct {
	Services  map[string]Service
	Databases []string
}

func (s ServiceYml) UnmarshalFromFile(path string) error {
	var err error
	_, err = os.Stat(path)
	service_yml.CheckError(err)

	yamlFile, err := ioutil.ReadFile(path)

	serviceYml := ServiceYml{
		Services:  map[string]Service{},
		Databases: make([]string, 1),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &serviceYml); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (s ServiceYml) MarshalToFile(path string) error {
	file, err := yaml.Marshal(s)

	file = []byte("# Generated with <3 by Cloud66\n\n" + string(s))

	err = ioutil.WriteFile(path, file, 0644)
	service_yml.CheckError(err)

	return nil
}
