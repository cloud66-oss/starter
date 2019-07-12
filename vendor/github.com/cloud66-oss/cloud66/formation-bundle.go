package cloud66

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FormationBundle struct {
	Version         string                  `json:"version"`
	Metadata        *Metadata               `json:"metadata"`
	Uid             string                  `json:"uid"`
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
	Uid           string `json:"uid"`
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
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}

type BundleTransformation struct { // this is just a placeholder for now
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Selector string   `json:"selector"`
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}

func CreateFormationBundle(formation Formation, app string, configurations []string, configstore []string) *FormationBundle {
	bundle := &FormationBundle{
		Version: "1",
		Metadata: &Metadata{
			App:         app,
			Timestamp:   time.Now().UTC(),
			Annotations: make([]string, 0), //just a placeholder before creating the real method
		},
		Uid:             formation.Uid,
		Name:            formation.Name,
		Tags:            formation.Tags,
		BaseTemplates:   createBaseTemplates(formation),
		Policies:        createPolicies(formation.Policies),
		Transformations: createTransformations(formation.Transformations),
		StencilGroups:   createStencilGroups(formation.StencilGroups),
		Configurations:  configurations,
		HelmReleases:    createHelmReleases(formation.HelmReleases),
		ConfigStore:     configstore,
	}
	return bundle
}

func createBaseTemplates(formation Formation) []*BundleBaseTemplates {
	baseTemplates := make([]*BundleBaseTemplates, 0)
	for _, stencil := range formation.Stencils {
		index := findIndexByRepoAndBranch(baseTemplates, stencil.BtrRepo, stencil.BtrBranch)
		if index == -1 {
			btrIndex := formation.FindIndexByRepoAndBranch(stencil.BtrRepo, stencil.BtrBranch)
			baseTemplates = append(baseTemplates, &BundleBaseTemplates{
				Name:     formation.BaseTemplates[btrIndex].Name,
				Repo:     formation.BaseTemplates[btrIndex].GitRepo,
				Branch:   formation.BaseTemplates[btrIndex].GitBranch,
				Stencils: createStencils(stencil),
			})
		} else {
			baseTemplates[index].Stencils = append(baseTemplates[index].Stencils, createStencil(stencil))
		}
	}
	return baseTemplates
}

func createStencils(stencil Stencil) []*BundleStencil {
	result := make([]*BundleStencil, 0)
	result = append(result, &BundleStencil{
		Uid:              stencil.Uid,
		Filename:         stencil.Filename,
		ContextID:        stencil.ContextID,
		TemplateFilename: stencil.TemplateFilename,
		Status:           stencil.Status,
		Tags:             stencil.Tags,
		Sequence:         stencil.Sequence,
	})

	return result
}

func createStencil(stencil Stencil) *BundleStencil {
	return &BundleStencil{
		Uid:              stencil.Uid,
		Filename:         stencil.Filename,
		ContextID:        stencil.ContextID,
		TemplateFilename: stencil.TemplateFilename,
		Status:           stencil.Status,
		Tags:             stencil.Tags,
		Sequence:         stencil.Sequence,
	}
}

func createStencilGroups(stencilGroups []StencilGroup) []*BundleStencilGroup {
	result := make([]*BundleStencilGroup, len(stencilGroups))
	for idx, st := range stencilGroups {
		result[idx] = &BundleStencilGroup{
			Name: st.Name,
			Uid:  st.Uid,
			Tags: st.Tags,
		}
	}

	return result
}

func createPolicies(policies []Policy) []*BundlePolicy {
	result := make([]*BundlePolicy, len(policies))
	for idx, st := range policies {
		result[idx] = &BundlePolicy{
			Uid:      st.Uid,
			Name:     st.Name,
			Selector: st.Selector,
			Sequence: st.Sequence,
			Tags:     st.Tags,
		}
	}

	return result
}

func createTransformations(transformations []Transformation) []*BundleTransformation {
	result := make([]*BundleTransformation, len(transformations))
	for idx, tr := range transformations {
		result[idx] = &BundleTransformation{
			Uid:      tr.Uid,
			Name:     tr.Name,
			Selector: tr.Selector,
			Sequence: tr.Sequence,
			Tags:     tr.Tags,
		}
	}

	return result
}

func (b *BundleStencil) AsStencil(bundlePath string) (*Stencil, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "stencils"), b.Filename)
	body, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return &Stencil{
		Uid:              b.Uid,
		Filename:         b.Filename,
		TemplateFilename: b.TemplateFilename,
		ContextID:        b.ContextID,
		Status:           b.Status,
		Tags:             b.Tags,
		Body:             string(body),
		Sequence:         b.Sequence,
	}, nil
}

func (b *BundlePolicy) AsPolicy(bundlePath string) (*Policy, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "policies"), b.Uid+".cop")
	body, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return &Policy{
		Uid:      b.Uid,
		Name:     b.Name,
		Selector: b.Selector,
		Sequence: b.Sequence,
		Body:     string(body),
		Tags:     b.Tags,
	}, nil
}

func (b *BundleTransformation) AsTransformation(bundlePath string) (*Transformation, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "transformations"), b.Uid+".js")
	body, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return &Transformation{
		Uid:      b.Uid,
		Name:     b.Name,
		Selector: b.Selector,
		Sequence: b.Sequence,
		Body:     string(body),
		Tags:     b.Tags,
	}, nil
}

func createHelmReleases(helmReleases []HelmRelease) []*BundleHelmRelease {
	result := make([]*BundleHelmRelease, len(helmReleases))
	for idx, hr := range helmReleases {
		filename := hr.DisplayName + "-values.yml"
		result[idx] = &BundleHelmRelease{
			ChartName:     hr.ChartName,
			DisplayName:   hr.DisplayName,
			Version:       hr.Version,
			RepositoryURL: hr.RepositoryURL,
			ValuesFile:    filename,
		}
	}

	return result
}

func (b *BundleHelmRelease) AsRelease(bundlePath string) (*HelmRelease, error) {
	var bodyString string = ""
	if b.ValuesFile != "" {
		filePath := filepath.Join(filepath.Join(bundlePath, "helm_releases"), b.ValuesFile)
		_, err := os.Stat(filePath)
		var body []byte
		if err != nil {
			body = nil
		} else {
			body, err = ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
		}
		bodyString = string(body)
	}

	return &HelmRelease{
		Uid:           b.Uid,
		ChartName:     b.ChartName,
		DisplayName:   b.DisplayName,
		RepositoryURL: b.RepositoryURL,
		Version:       b.Version,
		Body:          bodyString,
	}, nil
}

func (b *BundleStencilGroup) AsStencilGroup(bundlePath string) (*StencilGroup, error) {
	ext := ".json"
	body, err := ioutil.ReadFile(filepath.Join(bundlePath, "stencil_groups", b.Uid) + ext)
	if err != nil {
		return nil, err
	}

	return &StencilGroup{
		Uid:   b.Uid,
		Name:  b.Name,
		Tags:  b.Tags,
		Rules: string(body),
	}, nil
}

func findIndexByRepoAndBranch(base_templates []*BundleBaseTemplates, repo string, branch string) int {
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	for index, btr := range base_templates {
		if strings.TrimSpace(btr.Repo) == repo && strings.TrimSpace(btr.Branch) == branch {
			return index
		}
	}
	return -1
}
