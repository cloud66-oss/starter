package service_yml

import (
	"log"
	"strconv"

	"gopkg.in/yaml.v2"
)

func (t *TrafficMatches) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi *TrafficMatches
	err := unmarshal(multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		t = &TrafficMatches{single}
	} else {
		t = multi
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
		*p = append(*p, shortPortToLong(ports))
	case []interface{}:
		for _, vv := range ports {
			switch vv := vv.(type) {
			case string:
				*p = append(*p, shortPortToLong(vv))
			case int:
				vvi := strconv.Itoa(vv)
				*p = append(*p, shortPortToLong(vvi))
			case map[interface{}]interface{}:
				var longSyntaxStringPort tempPort
				var longSyntaxPort Port
				temp, er := yaml.Marshal(vv)
				CheckError(er)
				er = yaml.Unmarshal(temp, &longSyntaxStringPort)
				longSyntaxPort = stringToInt(longSyntaxStringPort)
				CheckError(er)
				*p = append(*p, longSyntaxPort)
			default:
				log.Fatal("Failed to unmarshal ")
			}
		}
	default:
		log.Fatal("Failed to unmarshal")

	}

	return nil
}
