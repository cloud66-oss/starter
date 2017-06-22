package transformer

import (
	"gopkg.in/yaml.v2"
	"log"
	"strconv"
)

func (e *BuildCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var build Build
	err := unmarshal(&build)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		e.BuildCommand = single
	} else {
		e.Build.Dockerfile = build.Dockerfile
	}
	return nil
}

func (sm *Volumes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		sm.Volumes = make([]string, 1)
		sm.Volumes[0] = single
	} else {
		sm.Volumes = multi
	}
	return nil
}

func (ef *EnvFile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		ef.EnvFile = make([]string, 1)
		ef.EnvFile[0] = single
	} else {
		ef.EnvFile = multi
	}
	return nil
}

func (sm *EnvVars) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m map[string]string
	var key, value string
	m = make(map[string]string, 1)
	err := unmarshal(&m)
	if err != nil {
		var multi []string
		err = unmarshal(&multi)
		if err != nil {
			var single string
			err = unmarshal(&single)
			if err != nil {
				return err
			}
			//get key, value and add to map m
			key, value = getKeyValue(single)
			m[key] = value
		}
		for i := 0; i < len(multi); i++ {
			key, value = getKeyValue(multi[i])
			m[key] = value
		}
		//get keys, values and add to map m
	}
	sm.EnvVars = m
	return nil
}

func (sm *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		sm.Command = make([]string, 1)
		sm.Command[0] = single
	} else {
		sm.Command = multi
	}
	return nil
}

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
				var longSyntaxPort Port
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
