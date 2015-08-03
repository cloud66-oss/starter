package webservers

import "github.com/cloud66/starter/packs"

type Unicorn struct {
	packs.WebServerBase
}

func (u *Unicorn) Names() []string {
	return []string{"unicorn"}
}

func (u *Unicorn) Port(command string) string {
	return u.WebServerBase.Port(u, command)
}

func (u *Unicorn) DefaultPort() string {
	return "8080"
}
