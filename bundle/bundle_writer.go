package bundle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud66-oss/starter/packs"
	"gopkg.in/go-yaml/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloud66-oss/cloud66"
	"github.com/cloud66-oss/starter/common"
	"github.com/sethvargo/go-password/password"
)

type ManifestBundle struct {
	Version         string                  `json:"version"`
	Metadata        *Metadata               `json:"metadata"`
	UID             string                  `json:"uid"`
	Name            string                  `json:"name"`
	StencilGroups   []*BundleStencilGroup   `json:"stencil_groups"`
	BaseTemplates   []*BundleBaseTemplates  `json:"base_templates"`
	Policies        []*BundlePolicy         `json:"policies"`
	Transformations []*BundleTransformation `json:"transformations"`
	Tags            []string                `json:"tags"`
	HelmReleases    []*BundleHelmRelease    `json:"helm_releases"`
	Configurations  []string                `json:"configuration"`
	ConfigStore     []string                `json:"configstore"`
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
	Name     string           `json:"name"`
	Repo     string           `json:"repo"`
	Branch   string           `json:"branch"`
	Stencils []*BundleStencil `json:"stencils"`
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

type BundleTransformation struct { // this is just a placeholder for now
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
	Stencils        []*StencilTemplate         `json:"stencils"`
	Policies        []*PolicyTemplate          `json:"policies"`
	Transformations []*TransformationsTemplate `json:"transformations"`
	HelmCharts      []*HelmChartTemplate       `json:"helm_charts"`
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
	Dependencies      []string `json:"dependencies"`
}

type PolicyTemplate struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}

type TransformationsTemplate struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}

type HelmChartTemplate struct {
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Tags               []string            `json:"tags"`
	ChartRepositoryUrl string              `json:"chart_repository_url"`
	ChartName          string              `json:"chart_name"`
	ChartVersion       string              `json:"chart_version"`
	Dependencies       []string            `json:"dependencies"`
	Modifiers          []*ModifierTemplate `json:"modifiers"`
}

type ModifierTemplate struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

func CreateSkycapFiles(outputDir string,
	templateRepository string,
	branch string,
	packName string,
	githubURL string,
	services []*common.Service,
	databases []common.Database,
	addGenericBtr bool) error {

	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}
	//Create .bundle directory structure if it doesn't exist
	tempFolder := os.TempDir()
	bundleFolder := filepath.Join(tempFolder, "bundle")
	defer os.RemoveAll(bundleFolder)
	err := CreateBundleFolderStructure(bundleFolder)
	if err != nil {
		return err
	}

	err = GenerateBundle(bundleFolder, templateRepository, branch, packName, githubURL, services, databases, false)
	if err != nil {
		return err
	}

	if addGenericBtr {
		err = GenerateBundle(bundleFolder, packs.GenericTemplateRepository(), branch, packs.GenericBundleSuffix(), packs.GithubURL(), services, databases, true)
		if err != nil {
			return err
		}
	}

	err = common.Tar(bundleFolder, filepath.Join(outputDir, "starter.bundle"))
	if err != nil {
		common.PrintError(err.Error())
	}
	fmt.Printf("Bundle is saved to starter.bundle\n")

	return err
}

func GenerateBundle(bundleFolder string,
	templateRepository string,
	branch string,
	packName string,
	githubURL string,
	services []*common.Service,
	databases []common.Database,
	isGenericBtr bool) error {
	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}

	//create manifest.json file and start filling
	manifestFile, err := loadManifest(bundleFolder)
	if err != nil {
		return err
	}

	templateJSON, err := generateTemplateJSONFromUpstreamFile(templateRepository, branch)
	if err != nil {
		return err
	}

	err = handleConfigStoreRecords(packName, services, databases, manifestFile, bundleFolder, isGenericBtr)
	if err != nil {
		return err
	}

	manifestFile, err = getRequiredStencils(
		templateJSON,
		templateRepository,
		branch,
		services,
		bundleFolder,
		manifestFile,
		githubURL)

	if err != nil {
		return err
	}

	if isGenericBtr {
		manifestFile, err = addDatabase(templateJSON, templateRepository, branch, bundleFolder, manifestFile, databases, githubURL)
		if err != nil {
			return err
		}
	} else {
		manifestFile, err = saveEnvVars(packName, getEnvVars(services, databases), manifestFile, bundleFolder)
		if err != nil {
			return err
		}
	}

	manifestFile, err = addPoliciesAndTransformations(manifestFile)

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
	return err
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

func getConfigStoreRecords(services []*common.Service, databases []common.Database, includeDatabases bool) ([]cloud66.BundledConfigStoreRecord, error) {
	result := make([]cloud66.BundledConfigStoreRecord, 0)
	if includeDatabases {
		for _, database := range databases {
			result = append(result, cloud66.BundledConfigStoreRecord{
				Scope: cloud66.BundledConfigStoreStackScope,
				ConfigStoreRecord: cloud66.ConfigStoreRecord{
					Key:      database.DockerImage + "." + "database",
					RawValue: base64.StdEncoding.EncodeToString([]byte("database")),
				},
			})

			generatedUsername, err := password.Generate(10, 5, 0, true, true)
			if err != nil {
				return nil, err
			}
			result = append(result, cloud66.BundledConfigStoreRecord{
				Scope: cloud66.BundledConfigStoreStackScope,
				ConfigStoreRecord: cloud66.ConfigStoreRecord{
					Key:      database.DockerImage + "." + "username",
					RawValue: base64.StdEncoding.EncodeToString([]byte(generatedUsername)),
				},
			})

			generatedPassword, err := password.Generate(64, 20, 0, false, true)
			if err != nil {
				return nil, err
			}
			result = append(result, cloud66.BundledConfigStoreRecord{
				Scope: cloud66.BundledConfigStoreStackScope,
				ConfigStoreRecord: cloud66.ConfigStoreRecord{
					Key:      database.DockerImage + "." + "password",
					RawValue: base64.StdEncoding.EncodeToString([]byte(generatedPassword)),
				},
			})

			result = append(result, cloud66.BundledConfigStoreRecord{
				Scope: cloud66.BundledConfigStoreStackScope,
				ConfigStoreRecord: cloud66.ConfigStoreRecord{
					Key:      database.DockerImage + "." + "host",
					RawValue: base64.StdEncoding.EncodeToString([]byte(database.DockerImage)),
				},
			})
		}
	}

	return result, nil
}

func setConfigStoreRecords(configStoreRecords []cloud66.BundledConfigStoreRecord, prefix string, manifestBundle *ManifestBundle, bundleFolder string) error {
	unmarshalledOutput := cloud66.BundledConfigStoreRecords{Records: configStoreRecords}
	marshalledOutput, err := yaml.Marshal(&unmarshalledOutput)
	if err != nil {
		return err
	}

	fileName := prefix + "-configstore.yml"
	filePath := filepath.Join(filepath.Join(bundleFolder, "configstore"), fileName)

	err = ioutil.WriteFile(filePath, marshalledOutput, 0600)
	if err != nil {
		return err
	}

	manifestBundle.ConfigStore = append(manifestBundle.ConfigStore, fileName)
	return nil
}

func CreateBundleFolderStructure(baseFolder string) error {
	var folders = []string{"stencils", "policies", "transformations", "stencil_groups", "helm_releases", "configurations", "configstore"}
	for _, subfolder := range folders {
		folder := filepath.Join(baseFolder, subfolder)
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRequiredStencils(
	templateJSON *TemplateJSON,
	templateRepository string,
	branch string,
	services []*common.Service,
	bundleFolder string,
	manifestFile *ManifestBundle,
	githubURL string) (*ManifestBundle, error) {

	initialComponentNames, err := getInitialComponentNames(templateJSON)
	if err != nil {
		return nil, err
	}
	requiredComponentNames, err := getRequiredComponentNames(templateJSON, initialComponentNames)
	if err != nil {
		return nil, err
	}

	var manifestStencils = make([]*BundleStencil, 0)
	requiredStencils := filterStencilsByRequiredComponentNames(templateJSON, requiredComponentNames)
	for _, stencil := range requiredStencils {
		if stencil.ContextType == "service" {
			for _, service := range services {
				manifestFile, manifestStencils, err = downloadAndAddStencil(
					service.Name,
					stencil,
					templateJSON.Name,
					manifestFile,
					bundleFolder,
					templateRepository,
					branch,
					manifestStencils,
				)
				if err != nil {
					return nil, err
				}
				// create entry in manifest file with formatted name
				// download and rename stencil file
			}
		} else {
			manifestFile, manifestStencils, err = downloadAndAddStencil(
				"",
				stencil,
				templateJSON.Name,
				manifestFile,
				bundleFolder,
				templateRepository,
				branch,
				manifestStencils,
			)
			if err != nil {
				return nil, err
			}
		}
	}
	var newTemplate BundleBaseTemplates
	newTemplate.Name = templateJSON.Name
	newTemplate.Repo = githubURL
	newTemplate.Branch = branch
	newTemplate.Stencils = manifestStencils

	manifestFile.BaseTemplates = append(manifestFile.BaseTemplates, &newTemplate)

	return manifestFile, nil
}

func loadManifest(bundleFolder string) (*ManifestBundle, error) {
	// TODO: if manifest file present, pick that up instead
	var manifest *ManifestBundle
	manifestPath := filepath.Join(bundleFolder, "manifest.json")
	if common.FileExists(manifestPath) {
		//open manifest.json file and cast it into the struct
		// open the template.json file and start downloading the stencils
		manifestFile, err := os.Open(manifestPath)
		if err != nil {
			return nil, err
		}
		manifestFileData, err := ioutil.ReadAll(manifestFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(manifestFileData, &manifest)
		if err != nil {
			return nil, err
		}
	} else {
		manifest = &ManifestBundle{
			Version:         "1",
			Metadata:        nil,
			UID:             "",
			Name:            "",
			StencilGroups:   make([]*BundleStencilGroup, 0),
			BaseTemplates:   make([]*BundleBaseTemplates, 0),
			Policies:        make([]*BundlePolicy, 0),
			Transformations: make([]*BundleTransformation, 0),
			Tags:            make([]string, 0),
			HelmReleases:    make([]*BundleHelmRelease, 0),
			Configurations:  make([]string, 0),
			ConfigStore:     make([]string, 0),
		}
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

func handleConfigStoreRecords(prefix string, services []*common.Service, databases []common.Database, manifestBundle *ManifestBundle, bundleFolder string, includeDatabases bool) error {
	configStoreRecords, err := getConfigStoreRecords(services, databases, includeDatabases)
	if err != nil {
		return err
	}

	if len(configStoreRecords) > 0 {
		err = setConfigStoreRecords(configStoreRecords, prefix, manifestBundle, bundleFolder)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadAndAddStencil(context string,
	stencil *StencilTemplate,
	btrShortname string,
	manifestFile *ManifestBundle,
	bundleFolder string,
	templateRepository string,
	branch string,
	manifestStencils []*BundleStencil) (*ManifestBundle, []*BundleStencil, error) {

	filename := ""
	if context != "" {
		filename = context + "_"
		if strings.HasPrefix(stencil.Filename, "_") {
			filename = "_" + filename + stencil.Filename[1:]
		}
	}
	filename = filename + stencil.Filename + "@" + btrShortname

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

func addPoliciesAndTransformations(manifestFile *ManifestBundle) (*ManifestBundle, error) {
	// manifestFile.Policies = make([]*BundlePolicy, 0)
	// manifestFile.Transformations = make([]*BundleTransformation, 0)
	return manifestFile, nil
}

func addDatabase(templateJSON *TemplateJSON, templateRepository, branch, bundleFolder string, manifestFile *ManifestBundle, databases []common.Database, githubURL string) (*ManifestBundle, error) {
	var helmReleases = manifestFile.HelmReleases
	for _, db := range databases {
		var release BundleHelmRelease

		switch db.Name {
		case "mysql":
			release.ChartName = db.Name
			release.DisplayName = db.Name
			release.Version = "1.2.1"
		case "postgresql":
			release.ChartName = db.Name
			release.DisplayName = db.Name
			release.Version = "5.3.11"
		case "redis":
			release.ChartName = db.Name
			release.DisplayName = db.Name
			release.Version = "8.0.14"
		case "mongodb":
			release.ChartName = db.Name
			release.DisplayName = db.Name
			release.Version = "5.20.3"
		default:
			common.PrintlnWarning("Database %s not supported\n", db.Name)
			continue
		}

		var applicableHelmChartTemplate *HelmChartTemplate
		for _, h := range templateJSON.Templates.HelmCharts {
			// TODO: maybe check the chart repository URL as well
			if h.ChartName == release.ChartName && h.ChartVersion == release.Version {
				applicableHelmChartTemplate = h
				break
			}
		}

		var valuesFile string
		if applicableHelmChartTemplate != nil {
			for _, modifier := range applicableHelmChartTemplate.Modifiers {
				if modifier.Type == "values.yml" {
					modifierContents, err := readStencilTemplateFile(templateRepository, branch, modifier.Filename)
					if err != nil {
						return nil, err
					}

					modifierBasename := path.Base(modifier.Filename)
					destinationFilename := filepath.Join(bundleFolder, "helm_releases", modifierBasename)
					err = ioutil.WriteFile(destinationFilename, modifierContents, 0644)
					if err != nil {
						return nil, err
					}

					valuesFile = modifierBasename
					break
				}
			}

			for _, dependency := range applicableHelmChartTemplate.Dependencies {
				temp := strings.SplitN(dependency, "/", 2)
				obj_type := temp[0]
				obj_name := temp[1]

				switch obj_type {
				case "stencils":
					stencilTemplate, err := getStencilTemplate(templateJSON, obj_name)
					if err != nil {
						return nil, err
					}

					baseTemplateRepoIndex, err := findIndexByRepoAndBranch(manifestFile.BaseTemplates, githubURL, branch)
					if err != nil {
						return nil, err
					}

					manifestFile, manifestFile.BaseTemplates[baseTemplateRepoIndex].Stencils, err = downloadAndAddStencil(
						"",
						stencilTemplate,
						templateJSON.Name,
						manifestFile,
						bundleFolder,
						templateRepository,
						branch,
						manifestFile.BaseTemplates[baseTemplateRepoIndex].Stencils)
					if err != nil {
						return nil, err
					}
				default:
					common.PrintlnWarning("Helm release dependency type %s not supported\n", obj_type)
					continue
				}
			}
		}

		release.UID = ""
		release.RepositoryURL = "https://kubernetes-charts.storage.googleapis.com/"
		release.ValuesFile = valuesFile
		helmReleases = append(helmReleases, &release)
	}
	manifestFile.HelmReleases = helmReleases
	return manifestFile, nil
}

func getStencilTemplate(templateJSON *TemplateJSON, stencil_name string) (*StencilTemplate, error) {
	for _, stencil := range templateJSON.Templates.Stencils {
		if stencil.Name == stencil_name {
			return stencil, nil
		}
	}
	return nil, errors.New("Stencil not found")
}

func findIndexByRepoAndBranch(base_templates []*BundleBaseTemplates, repo string, branch string) (int, error) {
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	for index, btr := range base_templates {
		if strings.TrimSpace(btr.Repo) == repo && strings.TrimSpace(btr.Branch) == branch {
			return index, nil
		}
	}
	return -1, errors.New("Base Template Repository not found inside the Bundle")
}

type color int

const (
	white color = 0
	grey  color = 1
	black color = 2
)

type DependencyInterface interface {
	getName() string
	getDependencies() []string
}

func (v StencilTemplate) getName() string {
	return v.Name
}

func (v StencilTemplate) getDependencies() []string {
	return v.Dependencies
}

func (v PolicyTemplate) getName() string {
	return v.Name
}

func (v PolicyTemplate) getDependencies() []string {
	return v.Dependencies
}

func (v TransformationsTemplate) getName() string {
	return v.Name
}

func (v TransformationsTemplate) getDependencies() []string {
	return v.Dependencies
}

func (v HelmChartTemplate) getName() string {
	return v.Name
}

func (v HelmChartTemplate) getDependencies() []string {
	return v.Dependencies
}

func getInitialComponentNames(templateJSON *TemplateJSON) ([]string, error) {
	result := make([]string, 0)
	for _, stencil := range templateJSON.Templates.Stencils {
		if stencil.MinUsage > 0 {
			fullyQualifiedStencilName, err := generateFullyQualifiedName(stencil)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedStencilName)
		}
	}
	return result, nil
}

func getRequiredComponentNames(templateJSON *TemplateJSON, initialComponentNames []string) ([]string, error) {
	// loop through them and get the full dependency tree
	requiredComponentNameMap := make(map[string]bool)
	for _, initialComponentName := range initialComponentNames {
		visited := make(map[string]color)

		err := getRequiredComponentNamesInternal(templateJSON, initialComponentName, initialComponentName, visited)
		if err != nil {
			return nil, err
		}

		for depencencyName, _ := range visited {
			requiredComponentNameMap[depencencyName] = true
		}
	}

	// get unique required component names
	requiredComponentNames := make([]string, 0)
	for requiredComponentName, _ := range requiredComponentNameMap {
		requiredComponentNames = append(requiredComponentNames, requiredComponentName)
	}
	return requiredComponentNames, nil
}

func getRequiredComponentNamesInternal(templateJSON *TemplateJSON, rootName string, name string, visited map[string]color) error {
	_, present := visited[name]
	if !present {
		visited[name] = white
	}

	currentColor, _ := visited[name]
	switch currentColor {
	case white:
		visited[name] = grey
	case grey:
		fmt.Printf("circular dependency for '%s' detected while processing dependency list of '%s'\n", name, rootName)
		return nil
	case black:
		return nil
	}

	templateDependencies, err := getTemplateDependencies(templateJSON, name)
	if err != nil {
		return err
	}
	for _, templateDependency := range templateDependencies {
		err := getRequiredComponentNamesInternal(templateJSON, rootName, templateDependency, visited)
		if err != nil {
			return err
		}
	}
	visited[name] = black

	return nil
}

func getTemplateDependencies(templateJSON *TemplateJSON, name string) ([]string, error) {
	nameParts := strings.Split(name, "/")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("dependency name '%s' should be 'TEMPLATE_TYPE/TEMPLATE_NAME', where TEMPLATE_TYPE is one of 'stencils', 'policies', 'transformations', or 'helm_charts'", name)
	}

	templateType := nameParts[0]
	templateName := nameParts[1]
	switch templateType {
	case "stencils":
		for _, v := range templateJSON.Templates.Stencils {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "policies":
		for _, v := range templateJSON.Templates.Policies {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "transformations":
		for _, v := range templateJSON.Templates.Transformations {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "helm_charts":
		for _, v := range templateJSON.Templates.HelmCharts {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	default:
		return nil, fmt.Errorf("dependency name '%s' should be 'TEMPLATE_TYPE/TEMPLATE_NAME', where TEMPLATE_TYPE is one of 'stencils', 'policies', 'transformations', or 'helm_charts'", name)
	}

	return nil, fmt.Errorf("could not find dependency with name '%s'", name)
}

func filterStencilsByRequiredComponentNames(templateJSON *TemplateJSON, requiredComponentNames []string) []*StencilTemplate {
	result := make([]*StencilTemplate, 0)
	for _, stencil := range templateJSON.Templates.Stencils {
		stencilRequired := false
		for _, requiredComponentName := range requiredComponentNames {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "stencils" && templateName == stencil.Name {
				stencilRequired = true
				break
			}
		}
		if stencilRequired {
			result = append(result, stencil)
		}
	}
	return result
}

func generateFullyQualifiedName(v DependencyInterface) (string, error) {
	name := v.getName()
	switch vt := v.(type) {
	case StencilTemplate, *StencilTemplate:
		return "stencils" + "/" + name, nil
	case PolicyTemplate, *PolicyTemplate:
		return "policies" + "/" + name, nil
	case TransformationsTemplate, *TransformationsTemplate:
		return "transformations" + "/" + name, nil
	case HelmChartTemplate, *HelmChartTemplate:
		return "helm_releases" + "/" + name, nil
	default:
		return "", fmt.Errorf("generateFullyQualifiedName missing definition for %T", vt)
	}
}

func generateTemplateJSONFromUpstreamFile(templateRepository, branch string) (*TemplateJSON, error) {
	templatesJSONData, err := readStencilTemplateFile(templateRepository, branch, "templates.json")
	if err != nil {
		return nil, err
	}

	var templateJSON TemplateJSON
	err = json.Unmarshal(templatesJSONData, &templateJSON)
	if err != nil {
		return nil, err
	}

	return &templateJSON, nil
}

func readStencilTemplateFile(templateRepository, branch, filename string) ([]byte, error) {
	temporaryFolder, err := ioutil.TempDir("", "bundle")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(temporaryFolder)

	//start download the template.json file
	downloadedFilePath, err := downloadStencilTemplateFile(templateRepository, branch, filename, temporaryFolder)
	if err != nil {
		fmt.Printf("Error while downloading file %s. The error is: %s\n", filename, err)
		return nil, err
	}

	// open the template.json file and start downloading the stencils
	downloadedFile, err := os.Open(downloadedFilePath)
	if err != nil {
		return nil, err
	}

	downloadedFileData, err := ioutil.ReadAll(downloadedFile)
	if err != nil {
		return nil, err
	}

	return downloadedFileData, nil
}

func downloadStencilTemplateFile(templateRepository, branch, filename, temporaryFolder string) (string, error) {
	manifestPath := templateRepository + filename // don't need to use filepath since it's a URL
	destinationFilename, err := common.GenerateRandomBase64String(32)
	if err != nil {
		return "", err
	}

	err = common.DownloadSingleFile(temporaryFolder, common.DownloadFile{URL: manifestPath, Name: destinationFilename}, branch)
	if err != nil {
		return "", err
	}

	return filepath.Join(temporaryFolder, destinationFilename), nil
}
