package packs

type WebServer interface {
	Names() []string
	Port(command *string) string
	ParsePort(command string) (hasFound bool, port string)
	RemovePortIfEnvVar(command string) string
	DefaultPort() string
}
