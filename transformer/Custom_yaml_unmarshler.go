package transformer

func (e *Build_Command) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var build Build
	err := unmarshal(&build)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		e.Build_Command = single
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

func (ef *Env_file) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		ef.Env_file = make([]string, 1)
		ef.Env_file[0] = single
	} else {
		ef.Env_file = multi
	}
	return nil
}

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

//unsupported warnings

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