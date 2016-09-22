package common

type EnvVar struct {
	Key   string
	Value string
}

func NewEnvMapping(key string, value string) *EnvVar {
	return &EnvVar{Key: key, Value: value}
}