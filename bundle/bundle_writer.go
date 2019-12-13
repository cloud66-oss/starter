package bundle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud66-oss/starter/bundle/bundles"
	"github.com/cloud66-oss/starter/bundle/templates"
	"log"
	"sort"

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

func CreateSkycapFiles(outputDir, templateRepository, branch, packName, githubURL string, services []*common.Service,
	databases []common.Database, addGenericBtr bool) error {

	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}

	//Create .bundle directory structure if it doesn't exist
	tempFolder := os.TempDir()
	bundleFolder := filepath.Join(tempFolder, "bundle")
	//bundleFolder := "/tmp/bundle"

	// cleanup the bundle folder
	defer func() {
		err := os.RemoveAll(bundleFolder)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := CreateBundleFolderStructure(bundleFolder)
	if err != nil {
		return err
	}

	err = GenerateBundleFiles(bundleFolder, templateRepository, branch, packName, githubURL, services, databases, false)
	if err != nil {
		return err
	}

	if addGenericBtr {
		err = GenerateBundleFiles(bundleFolder, packs.GenericTemplateRepository(), branch, packs.GenericBundleSuffix(), packs.GithubURL(), services, databases, true)
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

func GenerateBundleFiles(bundleFolder, templateRepository, branch, packName, githubURL string, services []*common.Service,
	databases []common.Database, isGenericBTR bool) error {

	if templateRepository == "" {
		//no stencil template defined for this pack, print an error and do nothing
		fmt.Printf("Sorry but there is no stencil template for this language/framework yet\n")
		return nil
	}

	// load the template file
	template, err := loadTemplate(templateRepository, branch)
	if err != nil {
		return err
	}

	// create bundle file to hold our structure
	bundle, err := loadBundle(bundleFolder)
	if err != nil {
		return err
	}

	// ensure our services are sorted (web first!)
	sortServices(services)

	// find components with min-usage 1
	minUsageComponents, err := getMinUsageComponents(template)
	if err != nil {
		return err
	}

	// find dependencies of the components above
	requiredComponents, err := getDependencyComponents(template, minUsageComponents)
	if err != nil {
		return err
	}

	err = handleConfigStoreRecords(packName, databases, bundle, bundleFolder, isGenericBTR)
	if err != nil {
		return err
	}

	// add stencils to the bundle
	bundle, err = addStencils(template, templateRepository, branch, services, bundleFolder, bundle, githubURL, requiredComponents)
	if err != nil {
		return err
	}

	if isGenericBTR {
		bundle, err = addDatabase(template, templateRepository, branch, bundleFolder, bundle, databases, githubURL)
		if err != nil {
			return err
		}
	} else {
		bundle, err = saveEnvVars(packName, getEnvVars(services), bundle, bundleFolder)
		if err != nil {
			return err
		}
	}

	err = addPolicies(bundle, template, templateRepository, branch, bundleFolder, requiredComponents)
	if err != nil {
		return err
	}

	err = addTransformations(bundle, template, templateRepository, branch, bundleFolder, requiredComponents)
	if err != nil {
		return err
	}

	err = addFilters(bundle, template, templateRepository, branch, bundleFolder, requiredComponents)
	if err != nil {
		return err
	}

	err = addWorkflows(bundle, template, templateRepository, branch, bundleFolder, isGenericBTR)
	if err != nil {
		return err
	}

	err = addMetadata(bundle)
	if err != nil {
		return err
	}

	err = saveManifest(bundleFolder, bundle)
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

func sortServices(services []*common.Service) {
	// sort services!
	sort.Slice(services[:], func(i, j int) bool {
		serviceName1 := strings.ToLower(services[i].Name)
		serviceName2 := strings.ToLower(services[j].Name)

		if serviceName1 == serviceName2 {
			// shouldn't happen?
			return false
		}

		// "web" should come out on top!
		if serviceName1 == "web" {
			return true
		}
		if serviceName2 == "web" {
			return false
		}

		// normal sort
		return serviceName1 < serviceName2
	})
}

func getEnvVars(services []*common.Service) map[string]string {
	var result = make(map[string]string)
	for _, envVarArray := range services {
		for _, envs := range envVarArray.EnvVars {
			result[envs.Key] = envs.Value
		}
	}
	return result
}

func getConfigStoreRecords(databases []common.Database, includeDatabases bool) ([]cloud66.BundledConfigStoreRecord, error) {
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

func setConfigStoreRecords(configStoreRecords []cloud66.BundledConfigStoreRecord, prefix string, bundle *bundles.Bundle, bundleFolder string) error {
	bundledConfigStoreRecords := cloud66.BundledConfigStoreRecords{Records: configStoreRecords}
	outputs, err := yaml.Marshal(&bundledConfigStoreRecords)
	if err != nil {
		return err
	}
	fileName := prefix + "-configstore.yml"
	filePath := filepath.Join(filepath.Join(bundleFolder, "configstore"), fileName)
	err = ioutil.WriteFile(filePath, outputs, 0600)
	if err != nil {

		return err
	}
	bundle.ConfigStore = append(bundle.ConfigStore, fileName)
	return nil
}

func CreateBundleFolderStructure(baseFolder string) error {
	var folders = []string{"stencils", "policies", "transformations", "stencil_groups", "helm_releases", "configurations", "configstore", "workflows", "filters"}
	for _, subFolder := range folders {
		folder := filepath.Join(baseFolder, subFolder)
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func addStencils(template *templates.Template, templateRepository string, branch string, services []*common.Service, bundleFolder string,
	bundle *bundles.Bundle, githubURL string, requiredComponents []string) (*bundles.Bundle, error) {

	var bundleStencils = make([]*bundles.Stencil, 0)
	templateStencils := filterStencilsByRequiredComponentNames(template, requiredComponents)

	stencilUsageMap := make(map[string]int)
	for _, templateStencil := range templateStencils {
		if templateStencil.Contextemplates == "service" {
			for _, service := range services {

				// check if we need this specific stencil
				if mustSkipStencil(templateStencil, service) {
					// this stencil isn't valid for this service
					continue
				}

				if stencilUsageMap[templateStencil.Filename] < templateStencil.MaxUsage {
					bundleStencil, err := downloadStencil(service.Name, templateStencil, template.Name, bundleFolder, templateRepository, branch)
					if err != nil {
						return nil, err
					}
					stencilUsageMap[templateStencil.Filename] += 1
					bundleStencils = append(bundleStencils, bundleStencil)
				} else {
					fmt.Printf("Skipping adding stencil '%s' for service '%s' because stencil max_usage exceeded\n", templateStencil.Name, service.Name)
				}
			}
		} else {
			if stencilUsageMap[templateStencil.Filename] < templateStencil.MaxUsage {
				bundleStencil, err := downloadStencil("", templateStencil, template.Name, bundleFolder, templateRepository, branch)
				if err != nil {
					return nil, err
				}
				stencilUsageMap[templateStencil.Filename] += 1
				bundleStencils = append(bundleStencils, bundleStencil)
			} else {
				fmt.Printf("Skipping adding stencil '%s' stencil because max_usage exceeded\n", templateStencil.Name)
			}
		}
	}
	var newTemplate bundles.BaseTemplate
	newTemplate.Name = template.Name
	newTemplate.Repo = githubURL
	newTemplate.Branch = branch
	newTemplate.Stencils = bundleStencils
	bundle.BaseTemplates = append(bundle.BaseTemplates, &newTemplate)
	return bundle, nil
}

func mustSkipStencil(templateStencil *templates.Stencil, service *common.Service) bool {
	// custom business logic for services in K8s
	mustSkip := false
	// if we have a service.yml template and the service we are dealing with
	// doesn't have external ports then we should ignore it
	if templateStencil.Filename == "service.yml" {
		mustSkip = true
		if service.Ports != nil || len(service.Ports) > 0 {
			for _, portMapping := range service.Ports {
				if portMapping.HTTP != "" || portMapping.HTTPS != "" || portMapping.TCP != "" || portMapping.UDP != "" {
					mustSkip = false
					continue
				}
			}
		}
	}
	return mustSkip
}

func addPolicies(bundle *bundles.Bundle, template *templates.Template, templateRepository, branch, bundleFolder string, requiredComponents []string) error {
	var bundlePolicies = bundle.Policies
	templatePolicys := filterPoliciesByRequiredComponentNames(template, requiredComponents)
	for _, templatePolicy := range templatePolicys {
		bundlePolicy, err := downloadPolicy(templatePolicy, template.Name, bundleFolder, templateRepository, branch)
		if err != nil {
			return err
		}
		bundlePolicies = append(bundlePolicies, bundlePolicy)
	}
	bundle.Policies = bundlePolicies
	return nil
}

func addTransformations(bundle *bundles.Bundle, template *templates.Template, templateRepository, branch, bundleFolder string, requiredComponents []string) error {
	var bundleTransformations = bundle.Transformations
	templateTransformations := filterTransformationsByRequiredComponentNames(template, requiredComponents)
	for _, templateTransformation := range templateTransformations {
		bundleTransformation, err := downloadTransformation(templateTransformation, template.Name, bundleFolder, templateRepository, branch)
		if err != nil {
			return err
		}
		bundleTransformations = append(bundleTransformations, bundleTransformation)
	}
	bundle.Transformations = bundleTransformations
	return nil
}

func addFilters(bundle *bundles.Bundle, template *templates.Template, templateRepository, branch, bundleFolder string, requiredComponents []string) error {
	var bundleFilters = bundle.Filters
	templateFilters := filterFiltersByRequiredComponentNames(template, requiredComponents)
	for _, templateFilter := range templateFilters {
		bundleFilter, err := downloadFilter(templateFilter, template.Name, bundleFolder, templateRepository, branch)
		if err != nil {
			return err
		}
		bundleFilters = append(bundleFilters, bundleFilter)
	}
	bundle.Filters = bundleFilters
	return nil
}

func addWorkflows(bundle *bundles.Bundle, template *templates.Template, templateRepository, branch, bundleFolder string, isGenericBTR bool) error {
	var bundleWorkflows = bundle.Workflows
	for _, templateWorkflow := range template.Templates.Workflows {
		bundleWorkflow, err := downloadWorkflow(templateWorkflow, template.Name, bundleFolder, templateRepository, branch, isGenericBTR)
		if err != nil {
			return err
		}
		bundleWorkflows = append(bundleWorkflows, bundleWorkflow)
	}
	bundle.Workflows = bundleWorkflows
	return nil
}

func loadBundle(bundleFolder string) (*bundles.Bundle, error) {
	var bundle *bundles.Bundle
	manifestPath := filepath.Join(bundleFolder, "manifest.json")
	if common.FileExists(manifestPath) {
		//open manifest.json file and cast it into the struct
		manifestFile, err := os.Open(manifestPath)
		if err != nil {
			return nil, err
		}
		manifestData, err := ioutil.ReadAll(manifestFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(manifestData, &bundle)
		if err != nil {
			return nil, err
		}
	} else {
		bundle = &bundles.Bundle{
			Version:         "1",
			Metadata:        nil,
			UID:             "",
			Name:            "",
			BaseTemplates:   make([]*bundles.BaseTemplate, 0),
			Policies:        make([]*bundles.Policy, 0),
			Transformations: make([]*bundles.Transformation, 0),
			Workflows:       make([]*bundles.Workflow, 0),
			HelmReleases:    make([]*bundles.HelmRelease, 0),
			Filters:         make([]*bundles.Filter, 0),
			Tags:            make([]string, 0),
			Configurations:  make([]string, 0),
			ConfigStore:     make([]string, 0),
		}
	}
	return bundle, nil
}

func saveManifest(bundleFolder string, content *bundles.Bundle) error {
	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(bundleFolder, "manifest.json")
	return ioutil.WriteFile(manifestPath, out, 0600)
}

func saveEnvVars(prefix string, envVars map[string]string, bundle *bundles.Bundle, bundleFolder string) (*bundles.Bundle, error) {
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
	var configs = bundle.Configurations
	bundle.Configurations = append(configs, filename)
	return bundle, nil
}

func handleConfigStoreRecords(prefix string, databases []common.Database, bundle *bundles.Bundle, bundleFolder string, includeDatabases bool) error {
	configStoreRecords, err := getConfigStoreRecords(databases, includeDatabases)
	if err != nil {
		return err
	}
	if len(configStoreRecords) > 0 {
		err = setConfigStoreRecords(configStoreRecords, prefix, bundle, bundleFolder)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadStencil(context string, templateStencil *templates.Stencil, btrShortName string, bundleFolder string, templateRepository string, branch string) (*bundles.Stencil, error) {
	filename := ""
	if context != "" {
		filename = context + "_"
		if strings.HasPrefix(templateStencil.Filename, "_") {
			filename = "_" + filename + templateStencil.Filename[1:]
		}
	}
	filename = filename + templateStencil.Filename
	filename, err := downloadComponent(templateStencil.Filename, filename, btrShortName, templateRepository, "stencils", bundleFolder, branch)
	if err != nil {
		return nil, err
	}

	// Add the entry to the manifest file
	var bundleStencil bundles.Stencil
	bundleStencil.UID = ""
	bundleStencil.Filename = filename
	bundleStencil.TemplateFilename = templateStencil.Filename
	bundleStencil.ContextID = context
	bundleStencil.Status = 2 // it means that the stencils still need to be deployed
	bundleStencil.Tags = templateStencil.Tags
	bundleStencil.Sequence = templateStencil.PreferredSequence

	return &bundleStencil, nil
}

func downloadPolicy(templatePolicy *templates.Policy, btrShortName, bundleFolder, templateRepository, branch string) (*bundles.Policy, error) {
	remoteFilename := templatePolicy.Filename
	localFilename := templatePolicy.Filename
	filename, err := downloadComponent(remoteFilename, localFilename, btrShortName, templateRepository, "policies", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundlePolicy := &bundles.Policy{
		UID:  "",
		Name: filename,
		Tags: templatePolicy.Tags,
	}
	return bundlePolicy, nil
}

func downloadTransformation(templateTransformation *templates.Transformation, btrShortName, bundleFolder, templateRepository, branch string) (*bundles.Transformation, error) {
	remoteFilename := templateTransformation.Filename
	localFilename := templateTransformation.Filename
	filename, err := downloadComponent(remoteFilename, localFilename, btrShortName, templateRepository, "transformations", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundleTransformation := &bundles.Transformation{
		UID:  "",
		Name: filename,
		Tags: templateTransformation.Tags,
	}
	return bundleTransformation, nil
}

func downloadFilter(templateFilter *templates.Filter, btrShortName, bundleFolder, templateRepository, branch string) (*bundles.Filter, error) {
	remoteFilename := templateFilter.Filename
	localFilename := templateFilter.Filename
	filename, err := downloadComponent(remoteFilename, localFilename, btrShortName, templateRepository, "filters", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundleFilter := &bundles.Filter{
		Name:        templateFilter.Name,
		Description: templateFilter.Description,
		Filename:    filename,
		Tags:        templateFilter.Tags,
	}
	return bundleFilter, nil
}

func downloadWorkflow(templateWorkflow *templates.Workflow, btrShortName, bundleFolder, templateRepository, branch string, isGenericBTR bool) (*bundles.Workflow, error) {
	remoteFilename := templateWorkflow.Filename
	localFilename := templateWorkflow.Filename
	filename, err := downloadComponent(remoteFilename, localFilename, btrShortName, templateRepository, "workflows", bundleFolder, branch)
	if err != nil {
		return nil, err
	}
	bundleWorkflow := &bundles.Workflow{
		Uid:     "",
		Name:    filename,
		Default: !isGenericBTR && templateWorkflow.Filename == "default.yml",
		Tags:    templateWorkflow.Tags,
	}
	return bundleWorkflow, nil
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

func addMetadata(bundle *bundles.Bundle) error {
	var metadata = &bundles.Metadata{
		Annotations: []string{"Generated by Cloud 66 starter"},
		App:         "starter",
		Timestamp:   time.Now().UTC(),
	}
	bundle.Metadata = metadata
	bundle.Name = "starter-formation"
	bundle.Tags = []string{"starter"}
	return nil
}

func addDatabase(template *templates.Template, templateRepository, branch, bundleFolder string, bundle *bundles.Bundle, databases []common.Database, githubURL string) (*bundles.Bundle, error) {
	var helmReleases = bundle.HelmReleases
	for _, db := range databases {
		var release bundles.HelmRelease

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

		var applicableHelmChartTemplate *templates.HelmRelease
		for _, h := range template.Templates.HelmCharts {
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
				objType := temp[0]
				objName := temp[1]

				switch objType {
				case "stencils":
					templateStencil, err := getStencilTemplate(template, objName)
					if err != nil {
						return nil, err
					}

					baseTemplateRepoIndex, err := findIndexByRepoAndBranch(bundle.BaseTemplates, githubURL, branch)
					if err != nil {
						return nil, err
					}

					bundleStencil, err := downloadStencil("", templateStencil, template.Name, bundleFolder, templateRepository, branch)
					if err != nil {
						return nil, err
					}
					bundle.BaseTemplates[baseTemplateRepoIndex].Stencils = append(bundle.BaseTemplates[baseTemplateRepoIndex].Stencils, bundleStencil)
				default:
					common.PrintlnWarning("Helm release dependency type %s not supported\n", objType)
					continue
				}
			}
		}

		release.UID = ""
		release.RepositoryURL = "https://kubernetes-charts.storage.googleapis.com/"
		release.ValuesFile = valuesFile
		helmReleases = append(helmReleases, &release)
	}
	bundle.HelmReleases = helmReleases
	return bundle, nil
}

func getStencilTemplate(template *templates.Template, stencilFilename string) (*templates.Stencil, error) {
	for _, stencil := range template.Templates.Stencils {
		if stencil.Filename == stencilFilename {
			return stencil, nil
		}
	}
	return nil, errors.New("stencil not found")
}

func findIndexByRepoAndBranch(baseTemplates []*bundles.BaseTemplate, repo string, branch string) (int, error) {
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	for index, btr := range baseTemplates {
		if strings.TrimSpace(btr.Repo) == repo && strings.TrimSpace(btr.Branch) == branch {
			return index, nil
		}
	}
	return -1, errors.New("base template repository not found inside the Bundle")
}

type color int

const (
	white color = 0
	grey  color = 1
	black color = 2
)

func getMinUsageComponents(template *templates.Template) ([]string, error) {
	result := make([]string, 0)
	// stencils
	for _, templateStencil := range template.Templates.Stencils {
		if templateStencil.MinUsage > 0 {
			fullyQualifiedStencilName, err := generateFullyQualifiedName(templateStencil)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedStencilName)
		}
	}
	// policies
	for _, templatePolicy := range template.Templates.Policies {
		if templatePolicy.MinUsage > 0 {
			fullyQualifiedPolicyName, err := generateFullyQualifiedName(templatePolicy)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedPolicyName)
		}
	}
	// transformations
	for _, templateTransformation := range template.Templates.Transformations {
		if templateTransformation.MinUsage > 0 {
			fullyQualifiedTransformationName, err := generateFullyQualifiedName(templateTransformation)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedTransformationName)
		}
	}
	// filters
	for _, templateFilter := range template.Templates.Filters {
		if templateFilter.MinUsage > 0 {
			fullyQualifiedFilterName, err := generateFullyQualifiedName(templateFilter)
			if err != nil {
				return nil, err
			}
			result = append(result, fullyQualifiedFilterName)
		}
	}
	return result, nil
}

func getDependencyComponents(template *templates.Template, initialComponentNames []string) ([]string, error) {
	// loop through them and get the full dependency tree
	requiredComponentNameMap := make(map[string]bool)
	for _, initialComponentName := range initialComponentNames {
		visited := make(map[string]color)
		err := getDependencyComponentsInternal(template, initialComponentName, initialComponentName, visited)
		if err != nil {
			return nil, err
		}
		for dependencyName := range visited {
			requiredComponentNameMap[dependencyName] = true
		}
	}
	// get unique required component names
	requiredComponents := make([]string, 0)
	for requiredComponentName := range requiredComponentNameMap {
		requiredComponents = append(requiredComponents, requiredComponentName)
	}
	return requiredComponents, nil
}

func getDependencyComponentsInternal(template *templates.Template, rootName string, name string, visited map[string]color) error {
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
	templateDependencies, err := getTemplateDependencies(template, name)
	if err != nil {
		return err
	}
	for _, templateDependency := range templateDependencies {
		err := getDependencyComponentsInternal(template, rootName, templateDependency, visited)
		if err != nil {
			return err
		}
	}
	visited[name] = black
	return nil
}

func getTemplateDependencies(template *templates.Template, name string) ([]string, error) {
	nameParts := strings.Split(name, "/")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("dependency name '%s' should be 'TEMPLATE_TYPE/TEMPLATE_NAME', where TEMPLATE_TYPE is one of 'stencils', 'policies', 'transformations', 'filters' or 'helm_charts'", name)
	}

	templateType := nameParts[0]
	templateName := nameParts[1]
	switch templateType {
	case "stencils":
		for _, v := range template.Templates.Stencils {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "policies":
		for _, v := range template.Templates.Policies {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "transformations":
		for _, v := range template.Templates.Transformations {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "helm_charts":
		for _, v := range template.Templates.HelmCharts {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	case "filters":
		for _, v := range template.Templates.Filters {
			if v.Name == templateName {
				return v.Dependencies, nil
			}
		}
	default:
		return nil, fmt.Errorf("dependency name '%s' should be 'TEMPLATE_TYPE/TEMPLATE_NAME', where TEMPLATE_TYPE is one of 'stencils', 'policies', 'transformations', or 'helm_charts'", name)
	}

	return nil, fmt.Errorf("could not find dependency with name '%s'", name)
}

func filterStencilsByRequiredComponentNames(template *templates.Template, requiredComponents []string) []*templates.Stencil {
	result := make([]*templates.Stencil, 0)
	for _, templateStencil := range template.Templates.Stencils {
		required := false
		for _, requiredComponentName := range requiredComponents {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "stencils" && templateName == templateStencil.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, templateStencil)
		}
	}
	return result
}

func filterPoliciesByRequiredComponentNames(template *templates.Template, requiredComponents []string) []*templates.Policy {
	result := make([]*templates.Policy, 0)
	for _, templatePolicy := range template.Templates.Policies {
		required := false
		for _, requiredComponentName := range requiredComponents {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "policies" && templateName == templatePolicy.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, templatePolicy)
		}
	}
	return result
}

func filterTransformationsByRequiredComponentNames(template *templates.Template, requiredComponents []string) []*templates.Transformation {
	result := make([]*templates.Transformation, 0)
	for _, templateTransformation := range template.Templates.Transformations {
		required := false
		for _, requiredComponentName := range requiredComponents {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "transformations" && templateName == templateTransformation.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, templateTransformation)
		}
	}
	return result
}

func filterFiltersByRequiredComponentNames(template *templates.Template, requiredComponents []string) []*templates.Filter {
	result := make([]*templates.Filter, 0)
	for _, templateFilter := range template.Templates.Filters {
		required := false
		for _, requiredComponentName := range requiredComponents {
			nameParts := strings.Split(requiredComponentName, "/")
			templateType := nameParts[0]
			templateName := nameParts[1]
			if templateType == "filters" && templateName == templateFilter.Name {
				required = true
				break
			}
		}
		if required {
			result = append(result, templateFilter)
		}
	}
	return result
}

func generateFullyQualifiedName(v templates.TemplateInterface) (string, error) {
	name := v.GetName()
	switch vt := v.(type) {
	case templates.Stencil, *templates.Stencil:
		return "stencils" + "/" + name, nil
	case templates.Policy, *templates.Policy:
		return "policies" + "/" + name, nil
	case templates.Transformation, *templates.Transformation:
		return "transformations" + "/" + name, nil
	case templates.HelmRelease, *templates.HelmRelease:
		return "helm_releases" + "/" + name, nil
	case templates.Workflow, *templates.Workflow:
		return "workflows" + "/" + name, nil
	case templates.Filter, *templates.Filter:
		return "filters" + "/" + name, nil
	default:
		return "", fmt.Errorf("generateFullyQualifiedName missing definition for %T", vt)
	}
}

func loadTemplate(templateRepository, branch string) (*templates.Template, error) {
	jsonData, err := readStencilTemplateFile(templateRepository, branch, "templates.json")
	if err != nil {
		return nil, err
	}
	var template templates.Template
	err = json.Unmarshal(jsonData, &template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func readStencilTemplateFile(templateRepository, branch, filename string) ([]byte, error) {
	temporaryFolder, err := ioutil.TempDir("", "bundle")
	if err != nil {
		return nil, err
	}

	// cleanup the bundle folder
	defer func() {
		err := os.RemoveAll(temporaryFolder)
		if err != nil {
			log.Fatal(err)
		}
	}()

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
