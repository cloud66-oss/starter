package templates

type Workflow struct {
	Name         string   `json:"name"`
	Filename     string   `json:"filename"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	Dependencies []string `json:"dependencies"`
}

func (v Workflow) GetName() string {
	return v.Name
}

func (v Workflow) GetDependencies() []string {
	return v.Dependencies
}
