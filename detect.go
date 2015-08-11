package main

import (
	"fmt"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/python"
	"github.com/cloud66/starter/packs/ruby"
)

func Detect(rootDir string) (packs.Pack, error) {
	node := node.Pack{}
	ruby := ruby.Pack{}
	python := python.Pack{}
	detectors := []packs.Detector{node.Detector(), ruby.Detector(), python.Detector()}
	var packs []packs.Pack

	for _, d := range detectors {
		if d.Detect(rootDir) {
			packs = append(packs, d.GetPack())
			common.PrintlnL0("Found %s application", d.GetPack().Name())
		}
	}

	if len(packs) == 0 {
		return nil, fmt.Errorf("Could not detect any of the supported frameworks")
	} else if len(packs) > 1 {
		return nil, fmt.Errorf("More than one framework detected")
	} else {
		return packs[0], nil
	}
}
