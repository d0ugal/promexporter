package main

import (
	"github.com/d0ugal/promexporter/linter/prometheuslabels"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(prometheuslabels.Analyzer)
}
