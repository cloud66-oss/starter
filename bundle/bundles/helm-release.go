package bundles

type HelmRelease struct {
	UID           string `json:"uid"`
	ChartName     string `json:"chart_name"`
	DisplayName   string `json:"display_name"`
	Version       string `json:"version"`
	RepositoryURL string `json:"repository_url"`
	ValuesFile    string `json:"values_file"`
}
