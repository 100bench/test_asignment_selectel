package loglinter_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	loglinter "github.com/100bench/loglinter"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), loglinter.Analyzer, "testcases")
}
