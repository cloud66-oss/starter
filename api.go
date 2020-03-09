package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/cloud66-oss/starter/bundle"
	"github.com/cloud66-oss/starter/common"
	service_yml "github.com/cloud66-oss/starter/definitions/service-yml"
	"github.com/cloud66-oss/starter/packs"
	"github.com/cloud66-oss/starter/packs/node"
	"github.com/cloud66-oss/starter/packs/php"
	"github.com/cloud66-oss/starter/packs/ruby"
	"github.com/cloud66-oss/starter/utils"
	"github.com/heroku/docker-registry-client/registry"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
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

var supportedBundlePacks = []packs.Pack{new(ruby.Pack)}

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
		&rest.Route{HttpMethod: "POST", PathExp: "/analyze/get-service", Func: a.getService},
		&rest.Route{HttpMethod: "POST", PathExp: "/analyze/get-bundle", Func: a.getBundle},
	)
	if err != nil {
		return err
	}

	api.SetApp(router)

	go func() {
		common.PrintL0("Starting API on %s\n", a.config.APIURL)
		common.PrintL1("API is now running...\n")

		packs := []packs.Pack{ /*new(compose_to_service_yml.Pack),*/ new(ruby.Pack), new(node.Pack), new(php.Pack) /*, new(service_yml_to_kubes.Pack)*/}

		for _, p := range packs {
			support := Language{}
			support.Name = p.Name()
			support.Files = p.FilesToBeAnalysed()

			if a.config.use_registry && p.Name() != "docker-compose" && p.Name() != "service.yml" {
				err := a.tryGetRemoteTags(p, support)
				if err != nil {
					//sleep then try again
					common.PrintL1("THROTTLE 60 seconds\n")
					time.Sleep(60 * time.Second)
					// try again
					err := a.tryGetRemoteTags(p, support)
					if err != nil {
						//sleep then try again ONE MORE TIME :/
						common.PrintL1("THROTTLE 60 seconds\n")
						time.Sleep(60 * time.Second)
						// try again
						err := a.tryGetRemoteTags(p, support)
						if err != nil {
							common.PrintError("Failed to start API due to: %s", err.Error())
							os.Exit(2)
						}
					}
				}
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

func (a *API) tryGetRemoteTags(p packs.Pack, support Language) error {
	url := "https://registry-1.docker.io/"
	username := "" // anonymous
	password := "" // anonymous
	hub, err := registry.New(url, username, password)
	if err != nil {
		return err
	}
	tags, err := hub.Tags("library/" + p.Name())
	if err != nil {
		return err
	}
	tags = Filter(tags, func(v string) bool {
		ok, _ := regexp.MatchString(`^\d+.\d+.\d+$`, v)
		return ok
	})
	p.SetSupportedLanguageVersions(tags)
	support.SupportedVersion = tags
	return nil
}

// routes system
func (a *API) ping(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("ok")
}

func (a *API) version(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(utils.Version)
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
			Version   string
			Framework string
			Packages  *common.Lister
		}{
			Version:   "latest",
			Framework: "express",
			Packages:  common.NewLister(),
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

	gitRepo := r.FormValue("git_repo")
	gitBranch := r.FormValue("git_branch")

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
	analysis := analyze_sourcecode(config, source_dir, "dockerfile,docker-compose,service,skycap", gitRepo, gitBranch)

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

func (a *API) getService(w rest.ResponseWriter, r *rest.Request) {
	uuid := uuid.NewV4().String()

	gitRepo := r.FormValue("git_repo")
	gitBranch := r.FormValue("git_branch")

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
	analysis := analyze_sourcecode(config, source_dir, "service", gitRepo, gitBranch)

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

func (a *API) getBundle(w rest.ResponseWriter, r *rest.Request) {
	uuid := uuid.NewV4().String()
	result := new(analysisResult)

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

	//analyze the service.yml file to get the bundle
	bundlePath, err := generate_bundle(source_dir)
	if err != nil {
		a.Error(w, err.Error(), 5, http.StatusInternalServerError)
		return
	}

	bundle, e := ioutil.ReadFile(bundlePath)
	if e != nil {
		a.Error(w, err.Error(), 6, http.StatusInternalServerError)
		return
	} else {
		result.SkycapBundle = base64.StdEncoding.EncodeToString(bundle)
	}

	// cleanup
	err = os.RemoveAll(source_dir)
	if err != nil {
		a.Error(w, err.Error(), 7, http.StatusInternalServerError)
		return
	}
	w.WriteJson(result)
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

func generate_bundle(sourcePath string) (string, error) {

	err := createBundleFromServiceFile(sourcePath, filepath.Join(sourcePath, "service.yml"), flagBTRBranch)

	if err != nil {
		return "", err
	}

	return filepath.Join(sourcePath, "starter.bundle"), nil
}

func analyze_sourcecode(config *Config, path string, generate string, gitRepo string, gitBranch string) *analysisResult {

	result, err := analyze(
		false,
		path,
		config.template_path,
		"production",
		true,
		true,
		generate,
		gitRepo,
		gitBranch,
		true)

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

		if strings.Contains(generate, "skycap") {
			bundle, e := ioutil.ReadFile(path + "/starter.bundle")
			if e != nil {
				// catch error
				result.SkycapBundle = ""
			} else {
				result.SkycapBundle = base64.StdEncoding.EncodeToString(bundle)
			}
		}
	}

	return result
}
func unzip(src, dest string) error {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}

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

func createBundleFromServiceFile(outputDir string,
	serviceFilePath string,
	branch string) error {

	serviceYml := service_yml.ServiceYml{}
	err := serviceYml.UnmarshalFromFile(filepath.Join(serviceFilePath))
	if err != nil {
		return err
	}

	var serviceContext packs.ServiceYAMLContextBase
	err = serviceContext.GenerateFromServiceYml(serviceYml)
	if err != nil {
		return err
	}

	services := serviceContext.Services
	databases := serviceContext.Dbs

	//Create .bundle directory structure if it doesn't exist
	tempFolder := os.TempDir()
	bundleFolder := filepath.Join(tempFolder, "bundle")
	defer os.RemoveAll(bundleFolder)
	err = bundle.CreateBundleFolderStructure(bundleFolder)
	if err != nil {
		return err
	}
	for _, pack := range getSupportedBundlePacks() {
		packServices := make([]*common.Service, 0)
		for _, service := range services {
			if contains(service.Tags, pack.FrameworkTag()) && contains(service.Tags, pack.LanguageTag()) {
				packServices = append(packServices, service)
			}
		}
		// generate the bundle for every supported tag
		err = bundle.GenerateBundleFiles(bundleFolder, pack.StencilRepositoryPath(), branch, pack.Name(), pack.PackGithubUrl(), packServices, databases, false)
		if err != nil {
			return err
		}
	}
	err = bundle.GenerateBundleFiles(bundleFolder, packs.GenericTemplateRepository(), branch, packs.GenericBundleSuffix(), packs.GithubURL(), services, databases, true)
	if err != nil {
		return err
	}

	err = common.Tar(bundleFolder, filepath.Join(outputDir, "starter.bundle"))
	if err != nil {
		common.PrintError(err.Error())
		return err
	}
	return err
}

func getSupportedBundlePacks() []packs.Pack {
	return supportedBundlePacks
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
