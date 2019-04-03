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
    HelmReleases   []*BundleHelmRelease     `json:"helm_releases"`
    Configurations []string                 `json:"configuration"`
}


type BundleHelmRelease struct {
    Name             string `json:"repo"`
    Version          string `json:"version"`
    RepositoryURL    string `json:"repository_url"`
	ValuesFile       string `json:"values_file"`
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
	manifestFile, err := loadManifest()
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

func getRequiredStencils(templateRepository string,
						branch string,
						outputDir string,
						services []*common.Service,
						skycapFolder string,
						manifestFile ManifestBundle,
						githubURL string) 	(ManifestBundle, error){

	//start download the template.json file
	tjPathfile, err := getStencilTemplateFile(templateRepository, skycapFolder, "templates.json", branch)
	if err != nil {
		fmt.Printf("Error while downloading the templates.json. err: %s", err)
		return ManifestBundle{}, err
	}
	// open the template.json file and start downloading the stencils
	templateJson, err := os.Open(tjPathfile)
	if err != nil {
		return ManifestBundle{}, err
	}

	templatesJsonData, err := ioutil.ReadAll(templateJson)
	if err != nil {
		return ManifestBundle{}, err
	}

	var templJson map[string]interface{}
	err = json.Unmarshal([]byte(templatesJsonData), &templJson)
	if err != nil {
		return ManifestBundle{}, err
	}

	var manifestStencils = make([]*BundleStencil, 0)

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
		}
	}
	var newTemplate BundleBaseTemplates
	newTemplate.Repo = githubURL
	newTemplate.Branch = branch
	newTemplate.Stencils = manifestStencils

	manifestFile.BaseTemplates =  append(manifestFile.BaseTemplates, &newTemplate)

	return manifestFile, nil
}


func loadManifest() (ManifestBundle, error) {
	var manifest ManifestBundle
	manifest.Version = "1"
	manifest.Metadata = nil
	manifest.Uid = ""
	manifest.Name = ""
	manifest.StencilGroups = make([]*BundleStencilGroup, 0)
	manifest.BaseTemplates = make([]*BundleBaseTemplates, 0)
	manifest.Policies = make([]*BundlePolicy, 0)
	manifest.Tags = make([]string, 0)
	manifest.HelmReleases = make([]*BundleHelmRelease, 0)
	manifest.Configurations = make([]string, 0)

	return manifest, nil
}

func saveManifest(skycapFolder string, content ManifestBundle) error {
	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(skycapFolder, "manifest.json")
	return ioutil.WriteFile(manifestPath, out, 0600)
}

func saveEnvVars(prefix string, envVars map[string]string , manifestFile ManifestBundle, skycapFolder string) (ManifestBundle, error) {
	filename := prefix+"-config"
	vars_path := filepath.Join(filepath.Join(skycapFolder,"configurations"), prefix+"-config")
	var fileOut string
	for key, value := range envVars {
		fileOut = fileOut+key+"="+value+"\n"
	}
	err := ioutil.WriteFile(vars_path, []byte(fileOut), 0600)
	if err!=nil {
		return ManifestBundle{}, err
	}
	var configs = manifestFile.Configurations
	manifestFile.Configurations = append(configs,filename)
	return manifestFile, nil
}

func downloadAndAddStencil(context string,
							stencil map[string]interface{},
							manifestFile ManifestBundle,
							skycapFolder string,
							templateRepository string,
							branch string,
							manifestStencils []*BundleStencil)	 (ManifestBundle, []*BundleStencil, error) {
	var filename = ""
	if context != "" {
		filename = context+"_"
	}
	filename = filename + stencil["filename"].(string)

	//download the stencil file
	stencil_path := templateRepository + stencil["filename"].(string)// don't need to use filepath since it's a URL
	stencils_folder := filepath.Join(skycapFolder, "stencils")
	down_err := common.DownloadSingleFile(stencils_folder, common.DownloadFile{URL: stencil_path, Name: filename}, branch)
	if down_err != nil {
		return ManifestBundle{}, nil, down_err
	}

	// Add the entry to the manifest file
	var tempStencil BundleStencil
	tempStencil.Uid= ""
	tempStencil.Filename = filename
	tempStencil.TemplateFilename = stencil["filename"].(string)
	tempStencil.ContextID = context
	tempStencil.Status = 2 // it means that the stencils still need to be deployed
	tempStencil.Tags = []string{"starter"}
	tempStencil.Sequence = int(stencil["preferred_sequence"].(float64))

	manifestStencils = append(manifestStencils, &tempStencil)

	return manifestFile, manifestStencils, nil
}

func addMetadatas(manifestFile ManifestBundle) (ManifestBundle, error) {
	var metadata Metadata
	metadata.Annotations = []string{"Generated by Cloud66 starter"}
	metadata.App = "starter"
	metadata.Timestamp = time.Now().UTC()
	manifestFile.Metadata = &metadata
	manifestFile.Name = "starter-formation"
	manifestFile.Tags = []string{"starter"}
	return manifestFile, nil
}

func addDatabase(manifestFile ManifestBundle, databases []common.Database) (ManifestBundle, error) {
	var helm_releases = make([]*BundleHelmRelease, 0)
	var release BundleHelmRelease
	for _, db := range databases {
		switch db.Name {
		case "mysql":
			release.Name = db.Name
			release.Version = "0.10.2"
		case "postgresql":
			release.Name = "postgresql"
			release.Version = "3.1.0"
		default:
			common.PrintlnWarning("Database %s not supported\n", db.Name)
			continue
		}
		release.RepositoryURL = "https://kubernetes-charts.storage.googleapis.com/"
		release.ValuesFile = ""
		helm_releases = append(helm_releases, &release)
	}
	manifestFile.HelmReleases = helm_releases
	return manifestFile, nil
}