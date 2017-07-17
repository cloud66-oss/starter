package service_yml

import "github.com/cloud66/starter/common"


func (sm *GitUrl) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"git_repo\" at the moment")
	return nil
}

func (sm *GitBranch) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"git_branch\" at the moment")
	return nil
}

func (sm *DockerfilePath) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"dockerfile_path\" at the moment")
	return nil
}

func (sm *Requires) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"requires\" at the moment")
	return nil
}

func (sm *BuildCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"build\" at the moment")
	return nil
}

func (sm *BuildRoot) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"build_root\" at the moment")
	return nil
}

func (sm *LogFolder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"log_folder\" at the moment")
	return nil
}

func (sm *DnsBehaviour) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"dns_behaviour\" at the moment")
	return nil
}

func (sm *UseHabitus) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"use_habitus\" at the moment")
	return nil
}

func (sm *UseHabitusStep) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"use_habitus_step\" at the moment")
	return nil
}

func (sm *Health) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"health\" at the moment")
	return nil
}

func (sm *PreStartSignal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"pre_start_signal\" at the moment")
	return nil
}

func (sm *PreStopSequence) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"pre_stop_sequence\" at the moment")
	return nil
}

func (sm *RestartOnDeploy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"restart_on_deploy\" at the moment")
	return nil
}

func (sm *TrafficMatches) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Kubernetes configuration format does not support \"traffic_matches\" at the moment")
	return nil
}

