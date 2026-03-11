package main

import (
	loglinter "github.com/100bench/loglinter"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

// AnalyzerPlugin is the entry point for golangci-lint plugin loading.
var AnalyzerPlugin analyzerPlugin

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{loglinter.Analyzer}
}

func main() {}
