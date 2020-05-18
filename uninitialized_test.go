package main

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata, err := filepath.Abs("./testdata")
	if err != nil {
		t.Fatalf("unable to resolve path to testdata: %v", err)
	}
	analysistest.Run(t, testdata, Analyzer, "./...")
}
