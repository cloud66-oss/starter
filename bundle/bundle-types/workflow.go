package bundle_types

type BundleWorkflow struct {
	Uid     string   `json:"uid"`
	Name    string   `json:"name"`
	Default bool     `json:"default"`
	Tags    []string `json:"tags"`
}
