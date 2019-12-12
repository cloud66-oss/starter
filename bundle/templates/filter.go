package templates

type Filter struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Filename     string   `json:"filename"`
	Tags         []string `json:"tags"`
	MinUsage     int      `json:"min_usage"`
	Dependencies []string `json:"dependencies"`
}

func (v Filter) GetName() string {
	return v.Name
}

func (v Filter) GetDependencies() []string {
	return v.Dependencies
}
