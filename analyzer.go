package main

import (
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/node"
	"github.com/cloud66/starter/packs/ruby"
)

func NewAnalyzer(packName string, rootDir string, environment string) packs.Analyzer {
	var a packs.Analyzer
	switch packName {
	case "ruby":
		a = &ruby.Analyzer{AnalyzerBase: packs.AnalyzerBase{RootDir: rootDir, Environment: environment}}
	case "node":
		a = &node.Analyzer{AnalyzerBase: packs.AnalyzerBase{RootDir: rootDir, Environment: environment}}
	}
	return a
}
