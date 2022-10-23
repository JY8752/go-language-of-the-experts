package main

import (
	"JY8752/go-language-of-the-experts/part1/complexity"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(complexity.Analyzer)
}
