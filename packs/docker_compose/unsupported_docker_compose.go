package docker_compose

import "github.com/cloud66/starter/common"

func (sm *CapAdd) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cap_add\" at the moment")
	return nil
}
func (sm *CapDrop) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cap_drop\" at the moment")
	return nil
}

func (sm *ContainerName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"container_name\" at the moment")
	return nil
}

func (sm *CgroupParent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cgroup_parent\" at the moment")
	return nil
}

func (sm *Devices) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"devices\" at the moment")
	return nil
}

func (sm *Dns) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"dns\" at the moment")
	return nil
}

func (sm *ExtraHosts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"hosts\" at the moment")
	return nil
}


func (sm *Isolation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"isolation\" at the moment")
	return nil
}


func (sm *Networks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"networks\" at the moment")
	return nil
}


func (sm *Secrets) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"secrets\" at the moment")
	return nil
}

func (sm *SecurityOpt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"security_opt\" at the moment")
	return nil
}

func (sm *UsernsMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"userns_mode\" at the moment")
	return nil
}

func (sm *Ulimits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"ulimits\" at the moment")
	return nil
}

func (sm *Healthcheck) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"healthcheck\" at the moment")
	return nil
}
func (sm *Logging) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"logging\" at the moment")
	return nil
}