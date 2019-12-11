package template_types

type JSON struct {
	Version     string         `json:"version"`
	Public      bool           `json:"public"`
	Name        string         `json:"name"`
	Icon        string         `json:"icon"`
	LongName    string         `json:"long_name"`
	Description string         `json:"description"`
	Templates   *TemplateTypes `json:"templates"`
}
