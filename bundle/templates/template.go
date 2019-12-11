package templates

type Template struct {
	Version         string            `json:"version"`
	Public          bool              `json:"public"`
	Name            string            `json:"name"`
	Icon            string            `json:"icon"`
	LongName        string            `json:"long_name"`
	Description     string            `json:"description"`
	Templates       *Templates        `json:"templates"`
}