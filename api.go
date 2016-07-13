package main

import (
	"net/http"
	"os"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/cloud66/starter/common"
	"io/ioutil"
)

// API holds starter API
type API struct {
	config *Config
}

type CodebaseAnalysis struct {
	Ok bool
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
	/* payload:
	path: path to the project to be examined
	generate: files to generate
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
		"",
		"",
		true,
		true,
		generate)

	if err != nil {
		a.handleError(err, w)
		return
	}

	


    analysis := CodebaseAnalysis{}
    analysis.Ok = result.OK
	analysis.Warnings = result.Warnings
    
    if result.OK {
        file, e := ioutil.ReadFile(path + "/Dockerfile")
	    if e != nil {
	    	// catch error
	    	analysis.Dockerfile = ""
	    } else {
    		analysis.Dockerfile = string(file)
    	}
    }
    w.WriteJson(analysis)
}

func (a *API) handleError(err error, w rest.ResponseWriter) {
	rest.Error(w, err.Error(), http.StatusBadRequest)
}
