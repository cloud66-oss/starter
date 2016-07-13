package main

import (
	"net/http"
	"os"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/cloud66/starter/common"
	"io/ioutil"
	"strings"
)

// API holds starter API
type API struct {
	config *Config
}

type CodebaseAnalysis struct {
	Ok bool
	Language string
	Framework string
    Warnings []string
    Dockerfile string
    Service string
    DockerCompose string
}

// NewAPI creates a new instance of the API
func NewAPI(configuration *Config) API {
	return API{config: configuration}
}

// StartAPI starts the API listeners
func (a *API) StartAPI() error {
	api := rest.NewApi()

	router, err := rest.MakeRouter(
		// system
		&rest.Route{HttpMethod: "GET", PathExp: "/ping", Func: a.ping},
		&rest.Route{HttpMethod: "GET", PathExp: "/version", Func: a.version},

		// parsing
		&rest.Route{HttpMethod: "POST", PathExp: "/analyze", Func: a.analyze},
	)
	if err != nil {
		return err
	}

	api.SetApp(router)

	go func() {
		common.PrintL0("Starting API on %s\n", a.config.APIURL)
		common.PrintL1("API is now running...\n")

		if err := http.ListenAndServe(a.config.APIURL, api.MakeHandler()); err != nil {
			common.PrintError("Failed to start API %s", err.Error())
			os.Exit(2)
		}
	}()

	return nil
}

// routes system
func (a *API) ping(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("ok")
}

func (a *API) version(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(VERSION)
}

// routes parsing
func (a *API) analyze(w rest.ResponseWriter, r *rest.Request) {
	/* 
	payload:
		{
			"path": "...", //path to the project to be examined
			"generate": "dockerfile,service,docker-compose" //files to generate 
		}
	
	response:
		{
		  "Ok": true,
		  "Language": "ruby",
		  "Warnings": null,
		  "Dockerfile": "...",
		  "Service": "...",
		  "DockerCompose": "..."
		}
	*/

	type payload struct {
		Path     string `json:"path"`
		Generate string `json:"generate"`
	}
	var request payload
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := request.Path
	generate := request.Generate

	result, err := analyze(
		false,
		path,
		config.template_path,
		"",
		true,
		true,
		generate)

	if err != nil {
		a.handleError(err, w)
		return
	}

	


    analysis := CodebaseAnalysis{}
    analysis.Language = result.Language
    analysis.Framework = result.Framework
    analysis.Ok = result.OK
	analysis.Warnings = result.Warnings
    
    if result.OK {
    	//always read the Dockerfile
        dockerfile, e := ioutil.ReadFile(path + "/Dockerfile")
	    if e != nil {
	    	// catch error
	    	analysis.Dockerfile = ""
	    } else {
    		analysis.Dockerfile = string(dockerfile)
    	}

    	if strings.Contains(generate, "service") {
	    	serviceymlfile, e := ioutil.ReadFile(path + "/service.yml")
		    if e != nil {
		    	// catch error
		    	analysis.Service = ""
		    } else {
	    		analysis.Service = string(serviceymlfile)
	    	}
    	}
    	if strings.Contains(generate, "docker-compose") {
	    	dockercomposeymlfile, e := ioutil.ReadFile(path + "/docker-compose.yml")
		    if e != nil {
		    	// catch error
		    	analysis.DockerCompose = ""
		    } else {
	    		analysis.DockerCompose = string(dockercomposeymlfile)
	    	}
    	}


    }
    w.WriteJson(analysis)
}

func (a *API) handleError(err error, w rest.ResponseWriter) {
	rest.Error(w, err.Error(), http.StatusBadRequest)
}
