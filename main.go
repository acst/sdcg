package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

var dir string

var ctx *hcl.EvalContext = &hcl.EvalContext{
	Variables: joinCtyVariables(stringGoTypeCty, stringConverterTypeCty),
	Functions: stringFuncTypeCty,
}

func init() {
	flag.StringVar(&dir, "dir", "", "directory to look for hcl files")
}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	files, err := parseHCLFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(files)
}

func parseFlags() error {
	flag.Parse()
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get working directory: %w", err)
		}
		dir = cwd
	}
	return nil
}

func parseHCLFiles() ([]*File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %w", err)
	}
	parser := hclparse.NewParser()
	var hclFiles []*File
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".hcl" {
			continue
		}
		hcl, diag := parser.ParseHCLFile(path.Join(dir, f.Name()))
		if diag.HasErrors() {
			return nil, fmt.Errorf("unable to parse %s: %w", f.Name(), diag)
		}
		var config File
		diag = gohcl.DecodeBody(hcl.Body, ctx, &config)
		if diag.HasErrors() {
			return nil, fmt.Errorf("unable to decode into go struct: %w", diag)
		}
		hclFiles = append(hclFiles, &config)
	}
	return hclFiles, nil
}
