package main

import "github.com/cloud66/starter/common"

type Docker_compose struct {
	Services map[string]Service
	Version string
}

type Service_yml struct {
	Services map[string]common.Service
	//Dbs []string
}

type Serviceyml struct{
	Services map[string]common.Service `yaml:"services,omitempty"`
	Dbs []string `yaml:"dbs,omitempty"`
}