package service_yml

type PreStopCommand struct {
	PreStopCommand []string `yaml:"pre_stop_command"`
}

type PostStartCommand struct{
	PostStartCommand []string `yaml:"post_start_command"`
}