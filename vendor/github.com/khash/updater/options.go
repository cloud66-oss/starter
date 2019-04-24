package updater

// Options defines options used in an Updater instance
type Options struct {
	RemoteURL            string
	VersionSpecsFilename string
	BinPattern           string
	Channel              string
	Silent               bool
}

// VersionSpecsURL returns the full URL for the VersionSpecs file
func (o *Options) VersionSpecsURL() string {
	return o.RemoteURL + o.VersionSpecsFilename
}

// BinURL returns the full URL pattern for the executable
func (o *Options) BinURL() string {
	return o.RemoteURL + o.BinPattern
}
