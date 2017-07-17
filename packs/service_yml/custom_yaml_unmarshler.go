package service_yml

import (
	"log"
	"strconv"
	"gopkg.in/yaml.v2"
)

func (p *Ports) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var ports interface{}

	err := unmarshal(&ports)
	if err != nil {
		log.Fatalf("yaml.Unmarshal: %v", err)
	}

	switch ports := ports.(type) {
	case string:
		p.ShortSyntax = append(p.ShortSyntax, ports)
	case []interface{}:
		for _, vv := range ports {
			switch vv := vv.(type) {
			case string:
				p.ShortSyntax = append(p.ShortSyntax, vv)
			case int:
				vvi := strconv.Itoa(vv)
				p.ShortSyntax = append(p.ShortSyntax, vvi)
			case map[interface{}]interface{}:
				var longSyntaxPort ServicePort
				temp, er := yaml.Marshal(vv)
				CheckError(er)
				er = yaml.Unmarshal(temp, &longSyntaxPort)
				CheckError(er)
				p.Port = append(p.Port, longSyntaxPort)
			default:
				log.Fatal("Failed to unmarshal ")
			}
		}
	default:
		log.Fatal("Failed to unmarshal")
	}

	return nil
}

func (com *PostStartCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var single string
	err := unmarshal(&single)
	if err!=nil{
		return err
	}
	com.PostStartCommand = make([]string,1)
	com.PostStartCommand[0] = single

	return nil
}

func (com *PreStopCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var single string
	err := unmarshal(&single)
	if err!=nil{
		return err
	}
	com.PreStopCommand = make([]string, 1)
	com.PreStopCommand[0] = single

	return nil
}

func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		c.Command = make([]string, 1)
		c.Command[0] = single
	} else {
		c.Command = multi
	}
	return nil
}

