package bundle

import (
	"encoding/json"
	"fmt"
	"github.com/cloud66-oss/starter/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type ManifestBundle struct {
    Version        string                   `json:"version"`
    Metadata       *Metadata                `json:"metadata"`
    Uid            string                   `json:"uid"`
    Name           string                   `json:"name"`
    StencilGroups  []*BundleStencilGroup    `json:"stencil_groups"`
    BaseTemplates  []*BundleBaseTemplates   `json:"base_template"`
    Policies       []*BundlePolicy          `json:"policies"`
    Tags           []string                 `json:"tags"`
    HelmReleases   []*BundleHelmReleases    `json:"helm_releases"`
    Configurations []string                 `json:"configuration"`
}


type BundleHelmReleases struct {
    Name             string `json:"repo"`
    Version          string `json:"version"`
    RepositoryURL    string `json:"repository_url"`
    Values           string `json:"values_file"`
}

type BundleConfiguration struct {
    Repo   string `json:"repo"`
    Branch string `json:"branch"`
}

type BundleBaseTemplates struct {
    Repo     string `json:"repo"`
    Branch   string `json:"branch"`
    Stencils []*BundleStencil `json:"stencils"`
}

type Metadata struct {
    App         string     `json:"app"`
    Timestamp   time.Time  `json:"timestamp"`
    Annotations []string   `json:"annotations"`
}

type BundleStencil struct {
    Uid              string   `json:"uid"`
    Filename         string   `json:"filename"`
    TemplateFilename string   `json:"template_filename"`
    ContextID        string   `json:"context_id"`
    Status           int      `json:"status"`
    Tags             []string `json:"tags"`
    Sequence         int      `json:"sequence"`
}

type BundleStencilGroup struct {
    Uid  string   `json:"uid"`
    Name string   `json:"name"`
    Tags []string `json:"tags"`
}

type BundlePolicy struct {
    Uid      string   `json:"uid"`
    Name     string   `json:"name"`
    Selector string   `json:"selector"`
    Tags     []string `json:"tags"`
}


func CreateSkycapFiles(outputDir string,
						templateDir string,
						templateRepository string,
						branch string,
						pack_name string,
						githubURL string,
						services []*common.Service,
						databases []common.Database) error {

	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}
	//Create .bundle directory structure if it doesn't exist
	skycapFolder := filepath.Join(outputDir, "skycap")
	err := createBundleFolderStructure(skycapFolder)
	defer os.RemoveAll(filepath.Join(skycapFolder, "temp"))

	//create manifest.json file and start filling
	manifestFile, err := loadManifest(skycapFolder, templateDir, pack_name)
	if err != nil {
		return err
	}

	manifestFile, err = saveEnvVars(pack_name, getEnvVars(services, databases), manifestFile, skycapFolder)
	if err != nil {
		return err
	}

	manifestFile, err = addDatabase(manifestFile, databases)

	manifestFile, err = getRequiredStencils(
		templateRepository,
		branch,
		outputDir,
		services,
		skycapFolder,
		manifestFile,
		githubURL)

	if err != nil {
		return err
	}

	manifestFile, err = addMetadatas(manifestFile)

	if err != nil {
		return err
	}

	saveManifest(skycapFolder, manifestFile)

	//fmt.Print(requiredStencilsList, err)


	err = saveManifest(skycapFolder, manifestFile)
	if err !=nil {
		return err
	}
	return err
}

// downloading templates from github and putting them into homedir
func getStencilTemplateFile(templateRepository string, skycapFolder string, filename string, branch string) (string, error) {

	//Download templates.json file
	manifest_path := templateRepository + filename // don't need to use filepath since it's a URL
	temp_folder := filepath.Join(skycapFolder, "temp")
	down_err := common.DownloadSingleFile(temp_folder, common.DownloadFile{URL: manifest_path, Name: filename}, branch)
	if down_err != nil {
		return "", down_err
	}
	return filepath.Join(temp_folder, filename), nil
}

func getEnvVars (servs []*common.Service, databases []common.Database) (map[string]string) {
	var envas = make(map[string]string)
	for _, envVarArray := range servs {
		for _, envs := range envVarArray.EnvVars {
			envas[envs.Key] = envs.Value
		}
	}
	return envas
}

func createBundleFolderStructure(baseFolder string) error {
	var folders = [6]string{"stencils", "policies", "stencil-group", "helm-releases", "temp", "configurations"}
	for _, subfolder := range folders {
		folder := filepath.Join(baseFolder, subfolder)
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRequiredStencils(templateRepository string, branch string, outputDir string, services []*common.Service,
	skycapFolder string, manifestFile map[string]interface{}, githubURL string) (map[string]interface{}, error){

	//start download the template.json file
	tjPathfile, err := getStencilTemplateFile(templateRepository, skycapFolder, "templates.json", branch)
	if err != nil {
		fmt.Printf("Error while downloading the templates.json. err: %s", err)
		return nil, err
	}
	// open the template.json file and start downloading the stencils
	templateJson, err := os.Open(tjPathfile)
	if err != nil {
		return nil, err
	}

	templatesJsonData, err := ioutil.ReadAll(templateJson)
	if err != nil {
		return nil, err
	}

	var templJson map[string]interface{}
	err = json.Unmarshal([]byte(templatesJsonData), &templJson)
	if err != nil {
		return nil, err
	}

	var manifestStencils = make([]interface{}, 0)
	//var stencilsArray []map[string]interface{}
	for i, data := range templJson {
		if i =="templates" {
			for _, stencils := range data.([]interface{}) {
				stencil := stencils.(map[string]interface{})
				if stencil["min_usage"].(float64) > 0 {
					if stencil["context_type"] == "service" {
						for _, service := range services {
							manifestFile, manifestStencils, err = downloadAndAddStencil(
								service.Name,
								stencil,
								manifestFile,
								skycapFolder,
								templateRepository,
								branch,
								manifestStencils)
							// create entry in manifest file with formatted name
							// download and rename stencil file
						}
					}else {
						manifestFile, manifestStencils, err = downloadAndAddStencil(
							"",
							stencil,
							manifestFile,
							skycapFolder,
							templateRepository,
							branch,
							manifestStencils)
					}
				}
			}
			// Do we need the db stencils?

		}
	}
	var templateMap = make(map[string]interface{},0)
	templateMap["repo"] = githubURL
	templateMap["branch"] = branch
	templateMap["stencils"] = manifestStencils

	manifestFile["base_templates"] =  append(manifestFile["base_templates"].([]interface{}), templateMap)

	return manifestFile, nil
}


func loadManifest(skycapFolder string, templateDir string, packName string) (map[string]interface{}, error) {

	//check local template file, if not present use the one in starter
	templateName :=  fmt.Sprintf("%s.bundle-manifest.json.template", packName)
	if !common.FileExists(filepath.Join(templateDir, templateName)) {
		templateName = "bundle-manifest.json.template" // fall back on generic template
	}
	// Open template file
	templates, err := os.Open(filepath.Join(templateDir, templateName))
	if err != nil {
		return nil, err
	}

	templatesData, err := ioutil.ReadAll(templates)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(templatesData), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func saveManifest(skycapFolder string, content map[string]interface{}) error {
	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(skycapFolder, "manifest.json")
	return ioutil.WriteFile(manifestPath, out, 0600)
}

func saveEnvVars(prefix string, envVars map[string]string , manifestFile map[string]interface{}, skycapFolder string) (map[string]interface{}, error) {
	filename := prefix+"-config"
	vars_path := filepath.Join(filepath.Join(skycapFolder,"configurations"), prefix+"-config")
	var fileOut string
	for key, value := range envVars {
		fileOut = fileOut+key+"="+value+"\n"
	}
	err := ioutil.WriteFile(vars_path, []byte(fileOut), 0600)
	if err!=nil {
		return nil, err
	}
	var configs = manifestFile["configurations"].([]interface {})
	manifestFile["configurations"] = append(configs,filename)
	return manifestFile, nil
}

func downloadAndAddStencil(context string, stencil map[string]interface{}, manifestFile map[string]interface{}, skycapFolder string, templateRepository string, branch string, manifestStencils []interface{}) (map[string]interface{}, []interface{}, error) {
	var filename = ""
	if context != "" {
		filename = context+"_"
	}
	filename = filename + stencil["filename"].(string)

	//download the stencil file
	stencil_path := templateRepository + stencil["filename"].(string) // don't need to use filepath since it's a URL
	stencils_folder := filepath.Join(skycapFolder, "stencils")
	down_err := common.DownloadSingleFile(stencils_folder, common.DownloadFile{URL: stencil_path, Name: filename}, branch)
	if down_err != nil {
		return nil, nil, down_err
	}

	// Add the entry to the manifest file
	var tempStencil = make(map[string]interface{}, 0)
	tempStencil["uid"] = ""
	tempStencil["filename"] = filename
	tempStencil["template_filename"] = stencil["filename"].(string)
	tempStencil["context_id"] = context
	tempStencil["status"] = 2 // it means that the stencils still need to be deployed
	var tags [1]string
	tags[0] = "starter"
	tempStencil["tags"] = tags
	tempStencil["sequence"] = stencil["preferred_sequence"]

	manifestStencils = append(manifestStencils, tempStencil)

	return manifestFile, manifestStencils, nil
}

func addMetadatas(manifestFile map[string]interface{}) (map[string]interface{}, error) {
	var details = make(map[string]interface{}, 0)
	details["info"] = "Generated by Cloud66 starter"

	var annotations = make([]map[string]interface{}, 0)
	annotations = append(annotations, details)

	var metadata = make(map[string]interface{}, 0)
	metadata["app"] = "starter"
	metadata["timestamp"] = time.Now().UTC()
	metadata["annotations"] = annotations

	manifestFile["metadata"] = metadata
	manifestFile["name"] = "starter-formation"

	var tags [1]string
	tags[0] = "starter"
	manifestFile["tags"] = tags

	return manifestFile, nil
}

func addDatabase(manifestFile  map[string]interface{}, databases []common.Database) (map[string]interface{}, error) {
	var helm_releases = make([]map[string]interface{}, 0)
	var release = make(map[string]interface{}, 0)
	for _, db := range databases {
		switch db.Name {
		case "mysql":
			release["name"] = db.Name
			release["version"] = "0.10.2"
		case "postgresql": //or postgresql ?
			release["name"] = "postgresql"
			release["version"] = "3.1.0"
		default:
			common.PrintlnWarning("Database %s not supported\n", db.Name)
			continue
		}
		release["repository_url"] = "https://kubernetes-charts.storage.googleapis.com/"
		release["values_file"] = ""
		helm_releases = append(helm_releases, release)
	}
	manifestFile["helm_releases"] = helm_releases
	return manifestFile, nil
}