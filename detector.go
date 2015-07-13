package main

import (
	"fmt"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/ruby"
)

func Detect(root string) (packs.Detector, error) {
	detectors := []packs.Detector{&ruby.Detector{}, &node.Detector{}}
	var found []packs.Detector

	for _, d := range detectors {
		if d.Detect(root) {
			found = append(found, d)
			fmt.Printf("%s Found %s application\n", common.MsgL0, d.Name())
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("Could not detect any of the supported frameworks")
	} else if len(found) > 1 {
		return nil, fmt.Errorf("More than one framework detected")
	} else {
		return found[0], nil
	}
}
