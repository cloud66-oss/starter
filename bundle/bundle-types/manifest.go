package bundle_types

type ManifestBundle struct {
	Version         string                  `json:"version"`
	Metadata        *Metadata               `json:"metadata"`
	UID             string                  `json:"uid"`
	Name            string                  `json:"name"`
	StencilGroups   []*StencilGroupBundle   `json:"stencil_groups"`
	BaseTemplates   []*BundleBaseTemplates  `json:"base_templates"`
	Policies        []*BundlePolicy         `json:"policies"`
	Transformations []*BundleTransformation `json:"transformations"`
	Workflows       []*BundleWorkflow       `json:"workflows"`
	Tags            []string                `json:"tags"`
	HelmReleases    []*BundleHelmRelease    `json:"helm_releases"`
	Configurations  []string                `json:"configuration"`
	ConfigStore     []string                `json:"configstore"`
}
