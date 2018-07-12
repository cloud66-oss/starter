package packs

import "github.com/cloud66-oss/starter/common"

type WebServerBase struct {
}

func (b *WebServerBase) Port(w WebServer, command *string) string {
	withoutPortEnvVar := w.RemovePortIfEnvVar(*command)
	*command = withoutPortEnvVar
	hasFound, port := w.ParsePort(*command)
	if hasFound {
		return port
	} else {
		return w.DefaultPort()
	}
}

func (w *WebServerBase) ParsePort(command string) (hasFound bool, port string) {
	return common.ParsePort(command)
}

func (w *WebServerBase) RemovePortIfEnvVar(command string) string {
	return common.RemovePortIfEnvVar(command)
}
