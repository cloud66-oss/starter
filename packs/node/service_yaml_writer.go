package node

import "github.com/cloud66/starter/packs"

type ServiceYAMLContext struct {
	packs.ServiceYAMLContextBase
}

type ServiceYAMLWriter struct {
	packs.ServiceYAMLWriterBase
}

func (w *ServiceYAMLWriter) Write(context *ServiceYAMLContext) error {
	return w.ServiceYAMLWriterBase.Write(context)
}
