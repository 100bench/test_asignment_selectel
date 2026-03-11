package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	loglinter "github.com/100bench/loglinter"
)

func main() {
	singlechecker.Main(loglinter.Analyzer)
}
