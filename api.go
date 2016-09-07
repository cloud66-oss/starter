package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/php"
	"github.com/cloud66/starter/packs/ruby"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// API holds starter API
type API struct {
	config *Config
}

type CodebaseAnalysis struct {
	Ok               bool
	Language         string
	Framework        string
	FrameworkVersion string
	Warnings         []string
	Dockerfile       string
	Service          string
	DockerCompose    string
}

type Language struct {
	Name  string
	Files []string
}

type SupportedLanguages struct {
	Languages []Language
}

type Dockerfile struct {
	Language string
	Base     string
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
		&rest.Route{HttpMethod: "GET", PathExp: "/analyze/supported", Func: a.supported},
		&rest.Route{HttpMethod: "GET", PathExp: "/analyze/dockerfiles", Func: a.dockerfiles},
		&rest.Route{HttpMethod: "POST", PathExp: "/analyze/upload", Func: a.upload},
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

func (a *API) dockerfiles(w rest.ResponseWriter, r *rest.Request) {
	packs := []packs.Pack{new(ruby.Pack), new(node.Pack), new(php.Pack)}
	dockerfiles := []Dockerfile{}
	for _, p := range packs {
		dockerfile := Dockerfile{}
		dockerfile.Language = p.Name()

		//parse base template
		templateName := fmt.Sprintf("%s.dockerfile.template", p.Name())
		tmpl, _ := template.ParseFiles(filepath.Join(config.template_path, templateName))

		var doc bytes.Buffer

		version := struct {
			Version  string
			Packages *common.Lister
		}{
			Version:  "latest",
			Packages: common.NewLister(),
		}

		err := tmpl.Execute(&doc, version)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dockerfile.Base = doc.String()

		dockerfiles = append(dockerfiles, dockerfile)
	}
	w.WriteJson(dockerfiles)
}

func (a *API) upload(w rest.ResponseWriter, r *rest.Request) {
	uuid := uuid.NewV4().String()

	 git_repo := r.FormValue("git_repo")
	 git_branch := r.FormValue("git_branch")

	 common.PrintL0("Param git_repo:  %s\n", git_repo)
	 common.PrintL0("Param git_branch: %s\n", git_branch)


	//save the file to a random location
	file, handler, err := r.FormFile("source")
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := "/tmp/" + uuid + "/" + handler.Filename
	source_dir := "/tmp/" + uuid
	err = os.MkdirAll(source_dir, 0777)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//unzip the file
	unzip(filename, source_dir)

	//analyse
	analysis := analyze_sourcecode(config, source_dir, "dockerfile,docker-compose,service", git_repo, git_branch)

	//cleanup
	err = os.RemoveAll(source_dir)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(analysis)
}

// routes parsing

func (a *API) supported(w rest.ResponseWriter, r *rest.Request) {
	/*
			response:
				{
		          "languages": [
		            {
		              "name": "ruby",
		              "files": [
		                "Gemfile",
		                "Procfile",
		                "config/database.yml"
		              ]
		            },
		             {
		              "name": "php",
		              "files": [
		                "composer.json"
		              ]
		            },
		             {
		              "name": "node",
		              "files": [
		                "package.json",
		                "Procfile"
		              ]
		            }
		          ]
		        }
	*/

	packs := []packs.Pack{new(ruby.Pack), new(node.Pack), new(php.Pack)}
	languages := SupportedLanguages{}
	for _, p := range packs {
		support := Language{}
		support.Name = p.Name()
		support.Files = p.FilesToBeAnalysed()
		languages.Languages = append(languages.Languages, support)
	}
	w.WriteJson(languages)
}

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
			  "Framework": "rails",
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
	analysis := analyze_sourcecode(config, path, generate)
	w.WriteJson(analysis)
}

func (a *API) handleError(err error, w rest.ResponseWriter) {
	rest.Error(w, err.Error(), http.StatusBadRequest)
}

func analyze_sourcecode(config *Config, path string, generate string, git_repo string, git_branch string) CodebaseAnalysis {
	analysis := CodebaseAnalysis{}
	result, err := analyze(
		false,
		path,
		config.template_path,
		"production",
		true,
		true,
		generate)

	if err != nil {
		common.PrintL0("%v", err.Error())
		return analysis
	}

	analysis.Language = result.Language
	analysis.Framework = result.Framework
	analysis.FrameworkVersion = result.FrameworkVersion
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
	return analysis
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
