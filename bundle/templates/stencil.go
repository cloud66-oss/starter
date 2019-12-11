package templates

type Stencil struct {
	Name              string   `json:"name"`
	FilenamePattern   string   `json:"filename_pattern"`
	Filename          string   `json:"filename"`
	Description       string   `json:"description"`
	Contextemplates       string   `json:"context_type"`
	Tags              []string `json:"tags"`
	PreferredSequence int      `json:"preferred_sequence"`
	Suggested         bool     `json:"suggested"`
	MinUsage          int      `json:"min_usage"`
	MaxUsage          int      `json:"max_usage"`
	Dependencies      []string `json:"dependencies"`
}

func (v Stencil) GetName() string {
	return v.Name
}

func (v Stencil) GetDependencies() []string {
	return v.Dependencies
}
