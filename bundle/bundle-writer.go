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
	Version        string                 `json:"version"`
	Metadata       *Metadata              `json:"metadata"`
	UID            string                 `json:"uid"`
	Name           string                 `json:"name"`
	StencilGroups  []*BundleStencilGroup  `json:"stencil_groups"`
	BaseTemplates  []*BundleBaseTemplates `json:"base_template"`
	Tags           []string               `json:"tags"`
	HelmReleases   []*BundleHelmRelease   `json:"helm_releases"`
	Configurations []string               `json:"configuration"`
}

type BundleHelmRelease struct {
	UID           string `json:"uid"`
	ChartName     string `json:"chart_name"`
	DisplayName   string `json:"display_name"`
	Version       string `json:"version"`
	RepositoryURL string `json:"repository_url"`
	ValuesFile    string `json:"values_file"`
}

type BundleBaseTemplates struct {
	Name         string                `json:"name"`
	Repo         string                `json:"repo"`
	Branch       string                `json:"branch"`
	Stencils     []*BundleStencil      `json:"stencils"`
	Policies     []*BundlePolicy       `json:"policies"`
	Transformers []*BundleTransformers `json:"transformers"`
}

type Metadata struct {
	App         string    `json:"app"`
	Timestamp   time.Time `json:"timestamp"`
	Annotations []string  `json:"annotations"`
}

type BundleStencil struct {
	UID              string   `json:"uid"`
	Filename         string   `json:"filename"`
	TemplateFilename string   `json:"template_filename"`
	ContextID        string   `json:"context_id"`
	Status           int      `json:"status"`
	Tags             []string `json:"tags"`
	Sequence         int      `json:"sequence"`
}

type BundleStencilGroup struct {
	UID  string   `json:"uid"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type BundlePolicy struct {
	UID      string   `json:"uid"`
	Name     string   `json:"name"`
	Selector string   `json:"selector"`
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}

type BundleTransformer struct { // this is just a placeholder for now
	UID  string   `json:"uid"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TemplateJSON struct {
	Version     string           `json:"version"`
	Public      bool             `json:"public"`
	Name        string           `json:"name"`
	Icon        string           `json:"icon"`
	LongName    string           `json:"long_name"`
	Description string           `json:"description"`
	Templates   *TemplatesStruct `json:"templates"`
}

type TemplatesStruct struct {
	Stencils     []*StencilTemplate      `json:"stencils"`
	Policies     []*PolicyTemplate       `json:"policies"`
	Transformers []*TransformersTemplate `json:"transformers"`
}

type StencilTemplate struct {
	Name              string   `json:"name"`
	FilenamePattern   string   `json:"filename_pattern"`
	Filename          string   `json:"filename"`
	Description       string   `json:"description"`
	ContextType       string   `json:"context_type"`
	Tags              []string `json:"tags"`
	PreferredSequence int      `json:"preferred_sequence"`
	Suggested         bool     `json:"suggested"`
	MinUsage          int      `json:"min_usage"`
	MaxUsage          int      `json:"max_usage"`
}

type PolicyTemplate struct{}

type TransformersTemplate struct{}

func CreateSkycapFiles(outputDir string,
	templateRepository string,
	branch string,
	packName string,
	githubURL string,
	services []*common.Service,
	databases []common.Database) error {

	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}
	//Create .bundle directory structure if it doesn't exist
	tempFolder := os.TempDir()
	bundleFolder := filepath.Join(tempFolder, "bundle")
	defer os.RemoveAll(bundleFolder)
	err := createBundleFolderStructure(bundleFolder)
	if err != nil {
		return err
	}
	//create manifest.json file and start filling
	manifestFile, err := loadManifest()
	if err != nil {
		return err
	}

	manifestFile, err = saveEnvVars(packName, getEnvVars(services, databases), manifestFile, bundleFolder)
	if err != nil {
		return err
	}

	manifestFile, err = addDatabase(manifestFile, databases)

	manifestFile, err = getRequiredStencils(
		templateRepository,
		branch,
		outputDir,
		services,
		bundleFolder,
		manifestFile,
		githubURL)

	if err != nil {
		return err
	}

	manifestFile, err = addMetadata(manifestFile)

	if err != nil {
		return err
	}

	err = saveManifest(bundleFolder, manifestFile)
	if err != nil {
		return err
	}

	// tarball
	err = os.RemoveAll(filepath.Join(bundleFolder, "temp"))
	if err != nil {
		common.PrintError(err.Error())
	}

	err = common.Tar(bundleFolder, filepath.Join(outputDir, "starter.bundle"))
	if err != nil {
		common.PrintError(err.Error())
	}
	fmt.Printf("Bundle is saved to starter.bundle\n")

	return err
}

// downloading templates from github and putting them into homedir
func getStencilTemplateFile(templateRepository string, tempFolder string, filename string, branch string) (string, error) {

	//Download templates.json file
	manifestPath := templateRepository + filename // don't need to use filepath since it's a URL
	downErr := common.DownloadSingleFile(tempFolder, common.DownloadFile{URL: manifestPath, Name: filename}, branch)
	if downErr != nil {
		return "", downErr
	}
	return filepath.Join(tempFolder, filename), nil
}

func getEnvVars(servs []*common.Service, databases []common.Database) map[string]string {
	var envas = make(map[string]string)
	for _, envVarArray := range servs {
		for _, envs := range envVarArray.EnvVars {
			envas[envs.Key] = envs.Value
		}
	}
	return envas
}

func createBundleFolderStructure(baseFolder string) error {
	var folders = [5]string{"stencils", "policies", "stencil-group", "helm-releases", "configurations"}
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
	bundleFolder string,
	manifestFile *ManifestBundle,
	githubURL string) (*ManifestBundle, error) {

	templateFolder := filepath.Join(os.TempDir(), "temp")
	err := os.MkdirAll(templateFolder, 0777)
	defer os.RemoveAll(templateFolder)
	if err != nil {
		return nil, err
	}
	//start download the template.json file
	tjPathfile, err := getStencilTemplateFile(templateRepository, templateFolder, "templates.json", branch)
	if err != nil {
		fmt.Printf("Error while downloading the templates.json. err: %s", err)
		return nil, err
	}
	// open the template.json file and start downloading the stencils
	templateJSON, err := os.Open(tjPathfile)
	if err != nil {
		return nil, err
	}

	templatesJSONData, err := ioutil.ReadAll(templateJSON)
	if err != nil {
		return nil, err
	}

	var templJSON TemplateJSON
	err = json.Unmarshal([]byte(templatesJSONData), &templJSON)
	if err != nil {
		return nil, err
	}

	var manifestStencils = make([]*BundleStencil, 0)
	for _, stencil := range templJSON.Templates.Stencils {
		if stencil.MinUsage > 0 {
			if stencil.ContextType == "service" {
				for _, service := range services {
					manifestFile, manifestStencils, err = downloadAndAddStencil(
						service.Name,
						stencil,
						manifestFile,
						bundleFolder,
						templateRepository,
						branch,
						manifestStencils)
					// create entry in manifest file with formatted name
					// download and rename stencil file
				}
			} else {
				manifestFile, manifestStencils, err = downloadAndAddStencil(
					"",
					stencil,
					manifestFile,
					bundleFolder,
					templateRepository,
					branch,
					manifestStencils)
			}
		}
	}
	var newTemplate BundleBaseTemplates
	newTemplate.Name = templJSON.Name
	newTemplate.Repo = githubURL
	newTemplate.Branch = branch
	newTemplate.Stencils = manifestStencils
	newTemplate.Policies = make([]*BundlePolicy, 0)
	newTemplate.Transformers = make([]*BundleTransformer, 0)

	manifestFile.BaseTemplates = append(manifestFile.BaseTemplates, &newTemplate)

	return manifestFile, nil
}

func loadManifest() (*ManifestBundle, error) {
	manifest := &ManifestBundle{
		Version:        "1",
		Metadata:       nil,
		UID:            "",
		Name:           "",
		StencilGroups:  make([]*BundleStencilGroup, 0),
		BaseTemplates:  make([]*BundleBaseTemplates, 0),
		Tags:           make([]string, 0),
		HelmReleases:   make([]*BundleHelmRelease, 0),
		Configurations: make([]string, 0),
	}
	return manifest, nil
}

func saveManifest(bundleFolder string, content *ManifestBundle) error {
	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(bundleFolder, "manifest.json")
	return ioutil.WriteFile(manifestPath, out, 0600)
}

func saveEnvVars(prefix string, envVars map[string]string, manifestFile *ManifestBundle, bundleFolder string) (*ManifestBundle, error) {
	filename := prefix + "-config"
	varsPath := filepath.Join(filepath.Join(bundleFolder, "configurations"), prefix+"-config")
	var fileOut string
	for key, value := range envVars {
		fileOut = fileOut + key + "=" + value + "\n"
	}
	err := ioutil.WriteFile(varsPath, []byte(fileOut), 0600)
	if err != nil {
		return nil, err
	}
	var configs = manifestFile.Configurations
	manifestFile.Configurations = append(configs, filename)
	return manifestFile, nil
}

func downloadAndAddStencil(context string,
	stencil *StencilTemplate,
	manifestFile *ManifestBundle,
	bundleFolder string,
	templateRepository string,
	branch string,
	manifestStencils []*BundleStencil) (*ManifestBundle, []*BundleStencil, error) {
	var filename = ""
	if context != "" {
		filename = context + "_"
	}
	filename = filename + stencil.Filename

	//download the stencil file
	stencilPath := templateRepository + "stencils/" + stencil.Filename // don't need to use filepath since it's a URL
	stencilsFolder := filepath.Join(bundleFolder, "stencils")
	downErr := common.DownloadSingleFile(stencilsFolder, common.DownloadFile{URL: stencilPath, Name: filename}, branch)
	if downErr != nil {
		return nil, nil, downErr
	}

	// Add the entry to the manifest file
	var tempStencil BundleStencil
	tempStencil.UID = ""
	tempStencil.Filename = filename
	tempStencil.TemplateFilename = stencil.Filename
	tempStencil.ContextID = context
	tempStencil.Status = 2 // it means that the stencils still need to be deployed
	tempStencil.Tags = []string{"starter"}
	tempStencil.Sequence = stencil.PreferredSequence

	manifestStencils = append(manifestStencils, &tempStencil)

	return manifestFile, manifestStencils, nil
}

func addMetadata(manifestFile *ManifestBundle) (*ManifestBundle, error) {
	var metadata = &Metadata{
		Annotations: []string{"Generated by Cloud 66 starter"},
		App:         "starter",
		Timestamp:   time.Now().UTC(),
	}
	manifestFile.Metadata = metadata
	manifestFile.Name = "starter-formation"
	manifestFile.Tags = []string{"starter"}
	return manifestFile, nil
}

func addDatabase(manifestFile *ManifestBundle, databases []common.Database) (*ManifestBundle, error) {
	var helmReleases = make([]*BundleHelmRelease, 0)
	var release BundleHelmRelease
	for _, db := range databases {
		switch db.Name {
		case "mysql":
			release.ChartName = db.Name
			release.DisplayName = "starter-" + db.Name
			release.Version = "0.10.2"
		case "postgresql":
			release.ChartName = db.Name
			release.DisplayName = "starter-" + db.Name
			release.Version = "3.1.0"
		default:
			common.PrintlnWarning("Database %s not supported\n", db.Name)
			continue
		}
		release.UID = ""
		release.RepositoryURL = "https://kubernetes-charts.storage.googleapis.com/"
		release.ValuesFile = ""
		helmReleases = append(helmReleases, &release)
	}
	manifestFile.HelmReleases = helmReleases
	return manifestFile, nil
}
