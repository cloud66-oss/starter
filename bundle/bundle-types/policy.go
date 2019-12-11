package bundle_types

type BundlePolicy struct {
	UID      string   `json:"uid"`
	Name     string   `json:"name"`
	Selector string   `json:"selector"`
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}
