package bundle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	btype "github.com/cloud66-oss/starter/bundle/bundle-types"
	ttype "github.com/cloud66-oss/starter/bundle/template-types"

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

	// find components with min-usage 1
	minUsageComponents, err := getMinUsageComponents(templateJSON)
	if err != nil {
		return err
	}
	//fmt.Println("MINUSAGE")
	//fmt.Println(minUsageComponents)

	// find dependencies of the components above
	requiredComponentNames, err := getDependencyComponents(templateJSON, minUsageComponents)
	if err != nil {
		return err
	}
	//fmt.Println("REQUIRED")
	//fmt.Println(requiredComponentNames)

	// add stencils to the bundle
	manifestFile, err = addStencils(templateJSON, templateRepository, branch, services, bundleFolder, manifestFile, githubURL, requiredComponentNames)
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

	manifestFile, err = addPolicies(templateJSON, templateRepository, branch, bundleFolder, manifestFile, requiredComponentNames)
	if err != nil {
		return err
	}

	manifestFile, err = addTransformations(templateJSON, templateRepository, branch, bundleFolder, manifestFile, requiredComponentNames)
	if err != nil {
		return err
	}

	manifestFile, err = addWorkflows(templateJSON, templateRepository, branch, bundleFolder, manifestFile, isGenericBtr)
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

			result = append(result, cloud66.BundledConfigStoreRecord{
				Scope: cloud66.BundledConfigStoreStackScope,
				ConfigStoreRecord: cloud66.ConfigStoreRecord{
					Key:      database.DockerImage + "." + "present",
					RawValue: base64.StdEncoding.EncodeToString([]byte("true")),
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

func setConfigStoreRecords(configStoreRecords []cloud66.BundledConfigStoreRecord, prefix string, manifestBundle *btype.ManifestBundle, bundleFolder string) error {
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
	var folders = []string{"stencils", "policies", "transformations", "stencil_groups", "helm_releases", "configurations", "configstore", "workflows"}
	for _, subfolder := range folders {
		folder := filepath.Join(baseFolder, subfolder)
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func addStencils(templateJSON *ttype.JSON, templateRepository string, branch string, services []*common.Service, bundleFolder string,
	manifestFile *btype.ManifestBundle, githubURL string, requiredComponentNames []string) (*btype.ManifestBundle, error) {

	var bundleStencils = make([]*btype.BundleStencil, 0)
	stencilTemplates := filterStencilsByRequiredComponentNames(templateJSON, requiredComponentNames)

	stencilUsageMap := make(map[string]int)
	for _, stencilTemplate := range stencilTemplates {
		if stencilTemplate.ContextType == "service" {
			for _, service := range services {
				if stencilUsageMap[stencilTemplate.Filename] < stencilTemplate.MaxUsage {
					bundleStencil, err := downloadStencil(service.Name, stencilTemplate, templateJSON.Name, manifestFile, bundleFolder, templateRepository, branch)
					if err != nil {
						return nil, err
					}
					stencilUsageMap[stencilTemplate.Filename] += 1
					bundleStencils = append(bundleStencils, bundleStencil)
					// create entry in manifest file with formatted name
					// download and rename stencil file
				} else {
					fmt.Printf("Skipping adding stencil '%s' for service '%s' because stencil max_usage exceeded\n", stencilTemplate.Name, service.Name)
				}
			}
		} else {
			if stencilUsageMap[stencilTemplate.Filename] < stencilTemplate.MaxUsage {
				bundleStencil, err := downloadStencil("", stencilTemplate, templateJSON.Name, manifestFile, bundleFolder, templateRepository, branch)
				if err != nil {
					return nil, err
				}
				stencilUsageMap[stencilTemplate.Filename] += 1
				bundleStencils = append(bundleStencils, bundleStencil)
			} else {
				fmt.Printf("Skipping adding stencil '%s' stencil because max_usage exceeded\n", stencilTemplate.Name)
			}
		}
	}
	var newTemplate btype.BundleBaseTemplates
	newTemplate.Name = templateJSON.Name
	newTemplate.Repo = githubURL
	newTemplate.Branch = branch
	newTemplate.Stencils = bundleStencils
	manifestFile.BaseTemplates = append(manifestFile.BaseTemplates, &newTemplate)
	return manifestFile, nil
}

func addPolicies(templateJSON *ttype.JSON, templateRepository string, branch string, bundleFolder string, manifestFile *btype.ManifestBundle, requiredComponentNames []string) (*btype.ManifestBundle, error) {
	var bundlePolicies = manifestFile.Policies
	policyTemplates := filterPoliciesByRequiredComponentNames(templateJSON, requiredComponentNames)
	for _, policyTemplate := range policyTemplates {
		bundlePolicy, err := downloadPolicy(policyTemplate, templateJSON.Name, bundleFolder, templateRepository, branch)
		if err != nil {
			return nil, err
		}
		bundlePolicies = append(bundlePolicies, bundlePolicy)
	}
	manifestFile.Policies = bundlePolicies
	return manifestFile, nil
}

func addTransformations(templateJSON *ttype.JSON, templateRepository string, branch string, bundleFolder string, manifestFile *btype.ManifestBundle, requiredComponentNames []string) (*btype.ManifestBundle, error) {
	var bundleTransformations = manifestFile.Transformations
	transformationTemplates := filterTransformationsByRequiredComponentNames(templateJSON, requiredComponentNames)
	for _, transformationTemplate := range transformationTemplates {
		bundleTransformation, err := downloadTransformation(transformationTemplate, templateJSON.Name, bundleFolder, templateRepository, branch)
		if err != nil {
			return nil, err
		}
		bundleTransformations = append(bundleTransformations, bundleTransformation)
	}
	manifestFile.Transformations = bundleTransformations
	return manifestFile, nil
}

func addWorkflows(templateJSON *ttype.JSON, templateRepository string, branch string, bundleFolder string, manifestFile *btype.ManifestBundle, isGenericBtr bool) (*btype.ManifestBundle, error) {
	var manifestWorkflows = manifestFile.Workflows
	var err error
	for _, workflow := range templateJSON.Templates.Workflows {
		manifestWorkflows, err = downloadAndAddWorkflow(
			workflow,
			templateJSON.Name,
			bundleFolder,
			templateRepository,
			branch,
			manifestWorkflows,
			isGenericBtr,
		)

		if err != nil {
			return nil, err
		}
	}
	manifestFile.Workflows = manifestWorkflows
	return manifestFile, nil
}

func loadManifest(bundleFolder string) (*btype.ManifestBundle, error) {
	// TODO: if manifest file present, pick that up instead
	var manifest *btype.ManifestBundle
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
		manifest = &btype.ManifestBundle{
			Version:         "1",
			Metadata:        nil,
			UID:             "",
			Name:            "",
			StencilGroups:   make([]*btype.StencilGroupBundle, 0),
			BaseTemplates:   make([]*btype.BundleBaseTemplates, 0),
			Policies:        make([]*btype.BundlePolicy, 0),
			Transformations: make([]*btype.BundleTransformation, 0),
			Workflows:       make([]*btype.BundleWorkflow, 0),
			Tags:            make([]string, 0),
			HelmReleases:    make([]*btype.BundleHelmRelease, 0),
			Configurations:  make([]string, 0),
			ConfigStore:     make([]string, 0),
		}
	}
	return manifest, nil
}

func saveManifest(bundleFolder string, content *btype.ManifestBundle) error {
	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(bundleFolder, "manifest.json")
	return ioutil.WriteFile(manifestPath, out, 0600)
}

func saveEnvVars(prefix string, envVars map[string]string, manifestFile *btype.ManifestBundle, bundleFolder string) (*btype.ManifestBundle, error) {
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

func handleConfigStoreRecords(prefix string, services []*common.Service, databases []common.Database, manifestBundle *btype.ManifestBundle, bundleFolder string, includeDatabases bool) error {
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

func downloadStencil(context string, stencilTemplate *ttype.Stencil, btrShortName string, manifestFile *btype.ManifestBundle, bundleFolder string, templateRepository string, branch string) (*btype.BundleStencil, error) {
	filename := ""
	if context != "" {
		filename = context + "_"
		if strings.HasPrefix(stencilTemplate.Filename, "_") {
			filename = "_" + filename + stencilTemplate.Filename[1:]
		}
	}
	filename = filename + stencilTemplate.Filename
	filename, err := downloadComponent(stencilTemplate.Filename, filename, btrShortName, templateRepository, "stencils", bundleFolder, branch)
	if err != nil {
		return nil, err
	}

	// Add the entry to the manifest file
	var bundleStencil btype.BundleStencil
	bundleStencil.UID = ""
	bundleStencil.Filename = filename
	bundleStencil.TemplateFilename = stencilTemplate.Filename
	bundleStencil.ContextID = context
	bundleStencil.Status = 2 // it means that the stencils still need to be deployed
	bundleStencil.Tags = stencilTemplate.Tags
	bundleStencil.Sequence = stencilTemplate.PreferredSequence

	return &bundleStencil, nil
}

func downloadPolicy(policyTemplate *ttype.Policy, btrShortName string, bundleFolder string, templateRepository string, branch string) (*btype.BundlePolicy, error) {
	filename := policyTemplate.Filename
	filename, err := downloadComponent(filename, filename, btrShortName, templateRepository, "policies", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundlePolicy := &btype.BundlePolicy{
		UID:  "",
		Name: filename,
		Tags: policyTemplate.Tags,
	}
	return bundlePolicy, nil
}

func downloadTransformation(transformationTemplate *ttype.Transformation, btrShortName string, bundleFolder string, templateRepository string, branch string) (*btype.BundleTransformation, error) {
	filename := transformationTemplate.Filename
	filename, err := downloadComponent(filename, filename, btrShortName, templateRepository, "transformations", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundleTransformation := &btype.BundleTransformation{
		UID:  "",
		Name: filename,
		Tags: transformationTemplate.Tags,
	}
	return bundleTransformation, nil
}

func downloadComponent(remoteFilename, localFilename string, btrShortName string, templateRepository string, componentName string, bundleFolder string, branch string) (string, error) {
	parts := strings.Split(localFilename, ".")
	if len(parts) > 1 {
		ext := parts[len(parts)-1]
		nameParts := parts[:len(parts)-1]
		name := strings.Join(nameParts[:], ".")
		localFilename = name + "@" + btrShortName + "." + ext
	} else {
		localFilename = localFilename + "@" + btrShortName
	}
	//download the file
	webPath := templateRepository + componentName + "/" + remoteFilename
	localFolder := filepath.Join(bundleFolder, componentName)
	downErr := common.DownloadSingleFile(localFolder, common.DownloadFile{URL: webPath, Name: localFilename}, branch)
	return localFilename, downErr
}

func downloadAndAddWorkflow(
	workflowTemplate *ttype.Workflow,
	btrShortname string,
	bundleFolder string,
	templateRepository string,
	branch string,
	manifestWorkflows []*btype.BundleWorkflow,
	isGenericBtr bool) ([]*btype.BundleWorkflow, error) {

	filename := workflowTemplate.Filename
	parts := strings.Split(filename, ".")

	if len(parts) > 1 {
		ext := parts[len(parts)-1]
		nameParts := parts[:len(parts)-1]
		name := strings.Join(nameParts[:], ".")
		filename = name + "@" + btrShortname + "." + ext
	} else {
		filename = filename + "@" + btrShortname
	}
	//download the stencil file
	workflowPath := templateRepository + "workflows/" + workflowTemplate.Filename // don't need to use filepath since it's a URL
	workflowsFolder := filepath.Join(bundleFolder, "workflows")
	downErr := common.DownloadSingleFile(workflowsFolder, common.DownloadFile{URL: workflowPath, Name: filename}, branch)

	if downErr != nil {
		return nil, downErr
	}

	var wk *btype.BundleWorkflow
	if isGenericBtr {
		wk = &btype.BundleWorkflow{
			Uid:     "",
			Name:    filename,
			Default: false,
			Tags:    workflowTemplate.Tags,
		}
	} else {
		wk = &btype.BundleWorkflow{
			Uid:     "",
			Name:    filename,
			Default: false,
			Tags:    workflowTemplate.Tags,
		}
		if workflowTemplate.Filename == "default.yml" {
			wk.Default = true
		}
	}
	manifestWorkflows = append(manifestWorkflows, wk)

	return manifestWorkflows, nil
}

func addMetadata(manifestFile *btype.ManifestBundle) (*btype.ManifestBundle, error) {
	var metadata = &btype.Metadata{
		Annotations: []string{"Generated by Cloud 66 starter"},
		App:         "starter",
		Timestamp:   time.Now().UTC(),
	}
	manifestFile.Metadata = metadata
	manifestFile.Name = "starter-formation"
	manifestFile.Tags = []string{"starter"}
	return manifestFile, nil
}

func addDatabase(templateJSON *ttype.JSON, templateRepository, branch, bundleFolder string, manifestFile *btype.ManifestBundle, databases []common.Database, githubURL string) (*btype.ManifestBundle, error) {
	var helmReleases = manifestFile.HelmReleases
	for _, db := range databases {
		var release btype.BundleHelmRelease

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

		var applicableHelmChartTemplate *ttype.HelmRelease
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

					bundleStencil, err := downloadStencil("", stencilTemplate, templateJSON.Name, manifestFile, bundleFolder, templateRepository, branch)
					if err != nil {
						return nil, err
					}
					manifestFile.BaseTemplates[baseTemplateRepoIndex].Stencils = append(manifestFile.BaseTemplates[baseTemplateRepoIndex].Stencils, bundleStencil)
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

func getStencilTemplate(templateJSON *ttype.JSON, stencil_filename string) (*ttype.Stencil, error) {
	for _, stencil := range templateJSON.Templates.Stencils {
		if stencil.Filename == stencil_filename {
			return stencil, nil
		}
	}
	return nil, errors.New("Stencil not found")
}

func findIndexByRepoAndBranch(base_templates []*btype.BundleBaseTemplates, repo string, branch string) (int, error) {
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

func getMinUsageComponents(templateJSON *ttype.JSON) ([]string, error) {
	result := make([]string, 0)
	// stencils
	for _, stencilTemplate := range templateJSON.Templates.Stencils {
		if stencilTemplate.MinUsage > 0 {
			fullyQualifiedStencilName, err := generateFullyQualifiedName(stencilTemplate)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedStencilName)
		}
	}
	// policies
	for _, policyTemplate := range templateJSON.Templates.Policies {
		if policyTemplate.MinUsage > 0 {
			fullyQualifiedPolicyName, err := generateFullyQualifiedName(policyTemplate)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedPolicyName)
		}
	}
	// transformations
	for _, transformationTemplate := range templateJSON.Templates.Transformations {
		if transformationTemplate.MinUsage > 0 {
			fullyQualifiedTransformationName, err := generateFullyQualifiedName(transformationTemplate)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedTransformationName)
		}
	}
	return result, nil
}

func getDependencyComponents(templateJSON *ttype.JSON, initialComponentNames []string) ([]string, error) {
	// loop through them and get the full dependency tree
	requiredComponentNameMap := make(map[string]bool)
	for _, initialComponentName := range initialComponentNames {
		visited := make(map[string]color)
		err := getRequiredComponentNamesInternal(templateJSON, initialComponentName, initialComponentName, visited)
		if err != nil {
			return nil, err
		}
		for dependencyName := range visited {
			requiredComponentNameMap[dependencyName] = true
		}
	}
	// get unique required component names
	requiredComponentNames := make([]string, 0)
	for requiredComponentName := range requiredComponentNameMap {
		requiredComponentNames = append(requiredComponentNames, requiredComponentName)
	}
	return requiredComponentNames, nil
}

func getRequiredComponentNamesInternal(templateJSON *ttype.JSON, rootName string, name string, visited map[string]color) error {
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

func getTemplateDependencies(templateJSON *ttype.JSON, name string) ([]string, error) {
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

func filterStencilsByRequiredComponentNames(templateJSON *ttype.JSON, requiredComponentNames []string) []*ttype.Stencil {
	result := make([]*ttype.Stencil, 0)
	for _, stencilTemplate := range templateJSON.Templates.Stencils {
		required := false
		for _, requiredComponentName := range requiredComponentNames {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "stencils" && templateName == stencilTemplate.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, stencilTemplate)
		}
	}
	return result
}

func filterPoliciesByRequiredComponentNames(templateJSON *ttype.JSON, requiredComponentNames []string) []*ttype.Policy {
	result := make([]*ttype.Policy, 0)
	for _, policyTemplate := range templateJSON.Templates.Policies {
		required := false
		for _, requiredComponentName := range requiredComponentNames {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "policies" && templateName == policyTemplate.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, policyTemplate)
		}
	}
	return result
}

func filterTransformationsByRequiredComponentNames(templateJSON *ttype.JSON, requiredComponentNames []string) []*ttype.Transformation {
	result := make([]*ttype.Transformation, 0)
	for _, transformationTemplate := range templateJSON.Templates.Transformations {
		required := false
		for _, requiredComponentName := range requiredComponentNames {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "transformations" && templateName == transformationTemplate.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, transformationTemplate)
		}
	}
	return result
}

func generateFullyQualifiedName(v ttype.TemplateInterface) (string, error) {
	name := v.GetName()
	switch vt := v.(type) {
	case ttype.Stencil, *ttype.Stencil:
		return "stencils" + "/" + name, nil
	case ttype.Policy, *ttype.Policy:
		return "policies" + "/" + name, nil
	case ttype.Transformation, *ttype.Transformation:
		return "transformations" + "/" + name, nil
	case ttype.HelmRelease, *ttype.HelmRelease:
		return "helm_releases" + "/" + name, nil
	case ttype.Workflow, *ttype.Workflow:
		return "workflows" + "/" + name, nil
	default:
		return "", fmt.Errorf("generateFullyQualifiedName missing definition for %T", vt)
	}
}

func generateTemplateJSONFromUpstreamFile(templateRepository, branch string) (*ttype.JSON, error) {
	templatesJSONData, err := readStencilTemplateFile(templateRepository, branch, "templates.json")
	if err != nil {
		return nil, err
	}

	var templateJSON ttype.JSON
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
