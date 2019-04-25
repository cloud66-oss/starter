package updater

import "fmt"

// VersionSpec is a single version descriptor
type VersionSpec struct {
	Version string `json:"version"`
	Channel string `json:"channel"`
	Force   bool   `json:"force"`
}

// VersionSpecs is a collection of VersionSpecs
type VersionSpecs struct {
	Versions []*VersionSpec `json:"versions"`
}

// GetVersionByChannel returns the applicable version by a channel
func (v *VersionSpecs) GetVersionByChannel(channel string) (*VersionSpec, error) {
	for idx, item := range v.Versions {
		if item.Channel == channel {
			return v.Versions[idx], nil
		}
	}

	return nil, fmt.Errorf("no version found for channel %s", channel)
}
