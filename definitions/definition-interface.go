package definitions

type Definition interface {
	UnmarshalFromFile(path string)
	MarshalToFile(path string)
}
