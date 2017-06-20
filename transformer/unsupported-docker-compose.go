package transformer

import "github.com/cloud66/starter/common"

type Links struct {
	Links []string `yaml:"links,omitempty"`
}

func (sm *Links) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"links\" at the moment")
	return nil
}

type Cap_add struct {
	Cap_add string
}

func (sm *Cap_add) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cap_add\" at the moment")
	return nil
}
type Cap_drop struct {
	Cap_drop string
}
func (sm *Cap_drop) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cap_drop\" at the moment")
	return nil
}

type Container_name struct {
	Container_name string
}

func (sm *Container_name) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"container_name\" at the moment")
	return nil
}

type Cgroup_parent struct {
	Cgroup_parent string
}

func (sm *Cgroup_parent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"cgroup_parent\" at the moment")
	return nil
}

type Devices struct {
	Devices string
}

func (sm *Devices) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"devices\" at the moment")
	return nil
}

type Dns struct {
	Dns string
}

func (sm *Dns) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"dns\" at the moment")
	return nil
}

type Extra_hosts struct {
	Extra_hosts []string
}

func (sm *Extra_hosts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"hosts\" at the moment")
	return nil
}

type Isolation struct {
	Isolation string
}

func (sm *Isolation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"isolation\" at the moment")
	return nil
}

type Networks struct {
	Networks []string
}

func (sm *Networks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"networks\" at the moment")
	return nil
}

type Secrets struct {
	Secrets string
}

func (sm *Secrets) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"secrets\" at the moment")
	return nil
}

type Security_opt struct {
	Security_opt string
}

func (sm *Security_opt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"security_opt\" at the moment")
	return nil
}

type Userns_mode struct {
	Userns_mode string
}

func (sm *Userns_mode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"userns_mode\" at the moment")
	return nil
}

type Ulimits struct {
	Ulimits string
}

func (sm *Ulimits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"ulimits\" at the moment")
	return nil
}

type Healthcheck struct {
	Healthcheck []string
}

func (sm *Healthcheck) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"healthcheck\" at the moment")
	return nil
}

type Logging struct {
	Logging []string
}

func (sm *Logging) UnmarshalYAML(unmarshal func(interface{}) error) error {
	common.PrintlnWarning("Service.yml format does not support \"logging\" at the moment")
	return nil
}