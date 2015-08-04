package packs

import "github.com/cloud66/starter/common"

type WebServerBase struct {
}

func (b *WebServerBase) Port(w WebServer, command string) string {
	hasFound, port := w.ParsePort(command)
	if hasFound {
		return port
	} else {
		return w.DefaultPort()
	}
}

func (w *WebServerBase) ParsePort(command string) (hasFound bool, port string) {
	return common.ParsePort(command)
}
