package docker_compose

import (
	"log"
	"strconv"
	"gopkg.in/yaml.v2"
)

func (b *Build) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single string
	err := unmarshal(&single)

	if err != nil {
		var bd BuildAux
		err := unmarshal(&bd)
		if err != nil {
			return err
		}
		b.Context = bd.Context
		b.Dockerfile = bd.Dockerfile
		b.Args = bd.Args
		b.CacheFrom = bd.CacheFrom
		b.Labels = bd.Labels
	}
	b.Context = single

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
		ef = &EnvFile{single}
	} else {
		*ef = multi
	}
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
		c = &Command{single}
	} else {
		*c = multi
	}
	return nil
}

func (c *Dns) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		c = &Dns{single}
	} else {
		*c = multi
	}
	return nil
}

func (c *DnsSearch) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		c = &DnsSearch{single}
	} else {
		*c = multi
	}
	return nil
}

func (c *Tmpfs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		c = &Tmpfs{single}
	} else {
		*c = multi
	}
	return nil
}

func (c *Entrypoint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		c = &Entrypoint{single}
	} else {
		*c = multi
	}
	return nil
}

func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var env_vars map[string]string
	var key, value string
	env_vars = make(map[string]string, 1)

	err := unmarshal(&env_vars)
	if err != nil {
		var multi []string
		err = unmarshal(&multi)
		if err != nil {
			return err
		}
		for i := 0; i < len(multi); i++ {
			key, value = getKeyValue(multi[i])
			env_vars[key] = value
		}
	}
	*e = env_vars
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
				var longSyntaxPort Port
				temp, er := yaml.Marshal(vv)
				CheckError(er)
				er = yaml.Unmarshal(temp, &longSyntaxPort)
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

func (s *Secrets) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var secrets interface{}

	err := unmarshal(&secrets)
	if err != nil {
		log.Fatalf("yaml.Unmarshal: %v", err)
	}

	switch secrets := secrets.(type) {
	case string:
		*s = append(*s, shortSecretToLong(secrets))
	case []interface{}:
		for _, vv := range secrets {
			switch vv := vv.(type) {
			case string:
				*s = append(*s, shortSecretToLong(vv))
			case map[interface{}]interface{}:
				var longSyntaxSecret Secret
				temp, er := yaml.Marshal(vv)
				CheckError(er)
				er = yaml.Unmarshal(temp, &longSyntaxSecret)
				CheckError(er)
				*s = append(*s, longSyntaxSecret)
			default:
				log.Fatal("Failed to unmarshal ")
			}
		}
	default:
		log.Fatal("Failed to unmarshal")
	}
	return nil
}


func (v *Volumes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var volumes interface{}
	err := unmarshal(&volumes)
	if err != nil {
		log.Fatalf("yaml.Unmarshal: %v", err)
	}

	switch volumes := volumes.(type) {
	case string:
		*v = append(*v, shortVolumeToLong(volumes))
	case []interface{}:
		for _, vv := range volumes {
			switch vv := vv.(type) {
			case string:
				*v = append(*v, shortVolumeToLong(vv))
			case map[interface{}]interface{}:
				var longSyntaxVolume Volume
				temp, er := yaml.Marshal(vv)
				CheckError(er)
				er = yaml.Unmarshal(temp, &longSyntaxVolume)
				CheckError(er)
				*v = append(*v, longSyntaxVolume)
			default:
				log.Fatal("Failed to unmarshal ")
			}
		}
	default:
		log.Fatal("Failed to unmarshal")
	}
	return nil
}

func (s *Limits) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var single int
	err := unmarshal(&single)
	if err != nil {
		var limit Limits
		err := unmarshal(&limit)
		if err != nil {
			return err
		}
		s.Hard = limit.Hard
		s.Soft = limit.Soft
	}else{
		s.Hard = single
		s.Soft = single
	}
	return nil
}
