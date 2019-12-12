package bundles

type Filter struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Filename    string   `json:"filename"`
	Tags        []string `json:"tags"`
}
