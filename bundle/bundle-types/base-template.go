package bundle_types

type BundleBaseTemplates struct {
	Name     string           `json:"name"`
	Repo     string           `json:"repo"`
	Branch   string           `json:"branch"`
	Stencils []*BundleStencil `json:"stencils"`
}
