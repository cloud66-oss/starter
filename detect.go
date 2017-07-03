package main

import (
	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/php"
	"github.com/cloud66/starter/packs/ruby"
	"github.com/cloud66/starter/packs/docker_compose"
	"fmt"
	"strings"
	"bufio"
	"os"
)

func Detect(rootDir string) ([]packs.Pack, error) {
	ruby := ruby.Pack{}
	node := node.Pack{}
	php := php.Pack{}
	dockercompose := docker_compose.Pack{}
	detectors := []packs.Detector{dockercompose.Detector(), ruby.Detector(), node.Detector(), php.Detector()}
	var packs []packs.Pack

	for _, d := range detectors {
		if d.Detect(rootDir) {
			packs = append(packs, d.GetPack())
			common.PrintlnL0("Found %s application", d.GetPack().Name())
		}
	}
	return packs, nil
}

func choosePack(detectedPacks []packs.Pack, noPrompt bool) (packs.Pack, error) {

	if len(detectedPacks) == 0 {

		return nil, fmt.Errorf("Could not detect any of the supported frameworks")

	} else if len(detectedPacks) > 1 {

		if noPrompt == false {

			common.PrintlnTitle("More than one framework detected. Please choose which of the following should be used:")

			for i := 0; i < len(detectedPacks); i++ {
				common.PrintlnTitle(strings.ToUpper(detectedPacks[i].Name()))
			}

			reader := bufio.NewReader(os.Stdin)
			var answer string
			answer, _ = reader.ReadString('\n')

			answer = strings.ToUpper(answer)

			for i := 0; i < len(detectedPacks); i++ {
				temp := strings.ToUpper(detectedPacks[i].Name()) + "\n"
				if answer == temp {
					return detectedPacks[i], nil
				}
			}
			return nil, fmt.Errorf("Starter was unable to match your answer")

		} else {

			common.PrintlnTitle("More than one framework detected! NoPrompt flag value is set to true.")

			for i := 0; i < len(detectedPacks); i++ {
				if detectedPacks[i].Name() == "docker-compose" {
					return detectedPacks[i], nil
				}
			}

			return nil, fmt.Errorf("Multiple frameworks detected. Unable to generate.")
		}
	} else {

		return detectedPacks[0], nil
	}

	return detectedPacks[0], nil
}
