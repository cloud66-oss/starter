package bundles

type Manifest struct {
	Version         string            `json:"version"`
	Metadata        *Metadata         `json:"metadata"`
	UID             string            `json:"uid"`
	Name            string            `json:"name"`
	StencilGroups   []*StencilGroup   `json:"stencil_groups"`
	BaseTemplates   []*BaseTemplate   `json:"base_templates"`
	Policies        []*Policy         `json:"policies"`
	Transformations []*Transformation `json:"transformations"`
	Workflows       []*Workflow       `json:"workflows"`
	Tags            []string          `json:"tags"`
	HelmReleases    []*HelmRelease    `json:"helm_releases"`
	Configurations  []string          `json:"configuration"`
	ConfigStore     []string          `json:"configstore"`
}
