package webservers

import "github.com/cloud66/starter/packs"

type Thin struct {
	packs.WebServerBase
}

func (t *Thin) Names() []string {
	return []string{"thin"}
}

func (t *Thin) Port(command *string) string {
	return t.WebServerBase.Port(t, command)
}

func (t *Thin) DefaultPort() string {
	return "3000"
}
