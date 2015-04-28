package packs

type Pack interface {
	Name() string
	PackVersion() string
	Detect() (bool, error)
	Compile() error
}
