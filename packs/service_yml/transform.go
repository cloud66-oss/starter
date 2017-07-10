package service_yml

import (
	"github.com/cloud66/starter/common"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
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
		Dbs:  []string{},
	}

	kubes := Kubes{
		Services: make(map[string]KubesService),
		Deployments: make(map[string]KubesDeployment),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &serviceYML); err!=nil{
		fmt.Println(err.Error())
	}

	kubes.Services, kubes.Deployments = copyToKubes(serviceYML, shouldPrompt, filename)

	file, err := yaml.Marshal(kubes)


	//Might need some formating for the "file"
	file = []byte("# Generated with <3 by Cloud66\n\n"+string(file))

	err = ioutil.WriteFile(formatTarget, file, 0644)
	if err!=nil{
		return err
	}

	//common.PrintlnTitle("In the kubes pack!")
	return nil
}

func copyToKubes (serviceYml ServiceYml, shouldPrompt bool, filepath string) (map[string]KubesService, map[string]KubesDeployment){



	return nil, nil
}
