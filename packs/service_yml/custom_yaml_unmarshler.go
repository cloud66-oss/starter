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

