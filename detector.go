package main

import (
	"fmt"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/ruby"
)

func Detect(rootDir string) (string, error) {
	detectors := []packs.Detector{&ruby.Detector{}, &node.Detector{}}
	var packs []string

	for _, d := range detectors {
		if d.Detect(rootDir) {
			packs = append(packs, d.PackName())
			fmt.Printf("%s Found %s application\n", common.MsgL0, d.PackName())
		}
	}

	if len(packs) == 0 {
		return "", fmt.Errorf("Could not detect any of the supported frameworks")
	} else if len(packs) > 1 {
		return "", fmt.Errorf("More than one framework detected")
	} else {
		return packs[0], nil
	}
}
