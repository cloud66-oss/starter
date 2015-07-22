package common

type PortMapping struct {
	Container string
	HTTP      string
	HTTPS     string
	TCP       string
	UDP       string
}

func NewPortMapping() *PortMapping {
	return &PortMapping{HTTP: "80", HTTPS: "443"}
}
