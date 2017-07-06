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
	"github.com/heroku/docker-registry-client/registry"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"regexp"
	"strconv"
	"github.com/cloud66/starter/packs/docker_compose"
)

// API holds starter API
type API struct {
	config *Config
}

type Language struct {
	Name             string
	Files            []string
	SupportedVersion []string
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

var languages = SupportedLanguages{}

func (a *API) Error(w rest.ResponseWriter, error string, error_code int, http_code int) {
	w.WriteHeader(http_code)
	err := w.WriteJson(map[string]string{"Error": error, "ErrorCode": strconv.Itoa(error_code)})
	if err != nil {
		panic(err)
	}
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

		packs := []packs.Pack{new(docker_compose.Pack), new(ruby.Pack), new(node.Pack), new(php.Pack)}
		for _, p := range packs {
			support := Language{}
			support.Name = p.Name()
			support.Files = p.FilesToBeAnalysed()

			if a.config.use_registry && p.Name()!="docker-compose"{
				url := "https://registry-1.docker.io/"
				username := "" // anonymous
				password := "" // anonymous
				hub, err := registry.New(url, username, password)
				if err != nil {
					common.PrintError("Failed to start API %s", err.Error())
					os.Exit(2)
				}
				tags, err := hub.Tags("library/" + p.Name())
				if err != nil {
					common.PrintError("Failed to start API %s", err.Error())
					os.Exit(2)
				}
				tags = Filter(tags, func(v string) bool {
					ok, _ := regexp.MatchString(`^\d+.\d+.\d+$`, v)
					return ok
				})
				p.SetSupportedLanguageVersions(tags)
				support.SupportedVersion = tags
			} else {
				support.SupportedVersion = p.GetSupportedLanguageVersions()
			}

			languages.Languages = append(languages.Languages, support)
		}

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
			Framework string
			Packages *common.Lister
		}{
			Version:  "latest",
			Framework: "express",
			Packages: common.NewLister(),
		}

		err := tmpl.Execute(&doc, version)
		if err != nil {
			a.Error(w, err.Error(), 1, http.StatusInternalServerError)
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

	//save the file to a random location
	file, handler, err := r.FormFile("source")
	if err != nil {
		a.Error(w, err.Error(), 2, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := "/tmp/" + uuid + "/" + handler.Filename
	source_dir := "/tmp/" + uuid
	err = os.MkdirAll(source_dir, 0777)

	if err != nil {
		a.Error(w, err.Error(), 3, http.StatusInternalServerError)
		return
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		a.Error(w, err.Error(), 4, http.StatusInternalServerError)
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
		a.Error(w, err.Error(), 5, http.StatusInternalServerError)
		return
	}

	if analysis != nil {
		w.WriteJson(analysis)
	} else {
		a.Error(w, "No supported language and/or framework detected", 6, http.StatusOK)
	}
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
		a.Error(w, err.Error(), 7, http.StatusInternalServerError)
		return
	}
	path := request.Path
	generate := request.Generate
	analysis := analyze_sourcecode(config, path, generate, "", "")
	if analysis != nil {
		w.WriteJson(analysis)
	} else {
		a.Error(w, "No supported language and/or framework detected", 6, http.StatusOK)
	}

}

func analyze_sourcecode(config *Config, path string, generate string, git_repo string, git_branch string) *analysisResult {

	result, err := analyze(
		false,
		path,
		config.template_path,
		"production",
		true,
		true,
		generate,
		git_repo,
		git_branch,
		false)

	if err != nil {
		common.PrintL0("%v", err.Error())
		return result
	}

	if result.Ok {
		//always read the Dockerfile
		dockerfile, e := ioutil.ReadFile(path + "/Dockerfile")
		if e != nil {
			// catch error
			result.Dockerfile = ""
		} else {
			result.Dockerfile = string(dockerfile)
		}

		if strings.Contains(generate, "service") {
			serviceymlfile, e := ioutil.ReadFile(path + "/service.yml")
			if e != nil {
				// catch error
				result.Service = ""
			} else {
				result.Service = string(serviceymlfile)
			}
		}
		if strings.Contains(generate, "docker-compose") {
			dockercomposeymlfile, e := ioutil.ReadFile(path + "/docker-compose.yml")
			if e != nil {
				// catch error
				result.DockerCompose = ""
			} else {
				result.DockerCompose = string(dockercomposeymlfile)
			}
		}

	}
	return result
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

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
