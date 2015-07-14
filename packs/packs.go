package packs

type Pack interface {
	PackVersion() string
	OutputFolder() string
	DefaultVersion() string
}
