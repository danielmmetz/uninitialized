package main

import (
	"github.com/danielmmetz/uninitialized/internal/uninitialized"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(uninitialized.NewAnalyzer())
}
