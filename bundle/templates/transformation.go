package templates

type Transformation struct {
	Name         string   `json:"name"`
	Filename     string   `json:"filename"`
	Dependencies []string `json:"dependencies"`
	MinUsage     int      `json:"min_usage"`
	MaxUsage     int      `json:"max_usage"`
	Tags         []string `json:"tags"`
}

func (v Transformation) GetName() string {
	return v.Name
}

func (v Transformation) GetDependencies() []string {
	return v.Dependencies
}
