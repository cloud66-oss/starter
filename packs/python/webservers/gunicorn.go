package webservers

import "github.com/cloud66-oss/starter/packs"

type Gunicorn struct {
	packs.WebServerBase
}

func (g *Gunicorn) Names() []string {
	return []string{"gunicorn"}
}

func (g *Gunicorn) Port(command *string) string {
	return g.WebServerBase.Port(g, command)
}

func (g *Gunicorn) DefaultPort() string {
	return "8000"
}
