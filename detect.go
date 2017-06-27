package main

import (
	"fmt"
	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/php"
	"github.com/cloud66/starter/packs/ruby"
	"github.com/cloud66/starter/packs/docker_compose"
	"bufio"
	"os"
	"strings"
)

func Detect(rootDir string) (packs.Pack, error) {
	ruby := ruby.Pack{}
	node := node.Pack{}
	php := php.Pack{}
	dockercompose := docker_compose.Pack{}
	detectors := []packs.Detector{ruby.Detector(), node.Detector(), php.Detector(), dockercompose.Detector()}
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
		common.PrintlnTitle("More than one framework detected. Please choose which of the following should be used:")
		for i:=0;i<len(packs);i++{
			common.PrintlnTitle(strings.ToUpper(packs[i].Name()))
		}

		reader := bufio.NewReader(os.Stdin)
		var answer string
		answer, _ = reader.ReadString('\n')

		answer = strings.ToUpper(answer)

		for i:=0;i<len(packs);i++ {
			temp := strings.ToUpper(packs[i].Name())+"\n"
			if answer == temp{
				return packs[i], nil
			}
		}

			return nil, fmt.Errorf("Starter was unable to match your answer")
	} else {
		return packs[0], nil
	}
}
