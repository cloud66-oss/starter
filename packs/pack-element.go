package packs

type PackElement struct {
	Pack Pack
}

func (e *PackElement) GetPack() Pack {
	return e.Pack
}
