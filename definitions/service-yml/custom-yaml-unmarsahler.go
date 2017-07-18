package service_yml

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
