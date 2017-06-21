package transformer

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

/*
func (sm *EnvVars) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		sm.EnvVars = make([]string, 1)
		sm.EnvVars[0] = single
	} else {
		sm.EnvVars = multi
	}
	return nil
}
*/

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
	var multiPorts []Port
	err := unmarshal(&multiPorts)
	if err != nil {
		var singlePort Port
		err := unmarshal(&singlePort)
		if err != nil {
			var multiString []string
			err := unmarshal(&multiString)
			if err != nil {
				var single string
				err := unmarshal(&single)
				if err != nil {
					return err
				}
				p.ShortSyntax = make([]string, 1)
				p.ShortSyntax[0] = single
			} else {
				p.ShortSyntax = multiString
			}
			return nil
		}
		p.Port = make([]Port, 1)
		p.Port[0] = singlePort
	} else {
		p.Port = multiPorts
	}
	return nil
}

/*

func (p *Ports) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var f interface{}
	err := unmarshal(&f)
	CheckError(err)

	m := f.(map[string]interface{})
	for k, v := range m {
		fmt.Println("the key is now %s", k)
		switch vv :=v.(type) {
		case Port:
			fmt.Println("It is a long syntax port!")
		case []Port:
			fmt.Println("It is a slice of longs")
		case []string:
			fmt.Println("It is a slice of short")
		case string:
			fmt.Println("It is a short syntax port with value %s!", vv)
		}

	}

	return nil
}
*/
