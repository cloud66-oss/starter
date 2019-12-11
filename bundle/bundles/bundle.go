package bundles

type Bundle struct {
	Version         string            `json:"version"`
	Metadata        *Metadata         `json:"metadata"`
	UID             string            `json:"uid"`
	Name            string            `json:"name"`
	BaseTemplates   []*BaseTemplate   `json:"base_templates"`
	Policies        []*Policy         `json:"policies"`
	Transformations []*Transformation `json:"transformations"`
	Workflows       []*Workflow       `json:"workflows"`
	HelmReleases    []*HelmRelease    `json:"helm_releases"`
	Filters         []*Filter         `json:"filters"`
	Tags            []string          `json:"tags"`
	Configurations  []string          `json:"configuration"`
	ConfigStore     []string          `json:"configstore"`
}
