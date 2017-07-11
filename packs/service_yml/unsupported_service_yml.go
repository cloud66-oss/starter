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
