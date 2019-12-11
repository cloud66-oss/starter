package bundles

type BaseTemplate struct {
	Name     string     `json:"name"`
	Repo     string     `json:"repo"`
	Branch   string     `json:"branch"`
	Stencils []*Stencil `json:"stencils"`
}
