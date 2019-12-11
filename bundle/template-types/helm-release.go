package template_types

type HelmRelease struct {
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Tags               []string    `json:"tags"`
	ChartRepositoryUrl string      `json:"chart_repository_url"`
	ChartName          string      `json:"chart_name"`
	ChartVersion       string      `json:"chart_version"`
	Dependencies       []string    `json:"dependencies"`
	Modifiers          []*Modifier `json:"modifiers"`
}

func (v HelmRelease) GetName() string {
	return v.Name
}

func (v HelmRelease) GetDependencies() []string {
	return v.Dependencies
}
