package packs

import "regexp"

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
	portPattern := regexp.MustCompile(`(?:-p|--port=)[[:blank:]]*(\d+)`)
	ports := portPattern.FindAllStringSubmatch(command, -1)
	if len(ports) != 1 {
		return false, ""
	} else {
		return true, ports[0][1]
	}
}
