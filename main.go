package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2"
)

var dir string

var hclCtx *hcl.EvalContext = &hcl.EvalContext{
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
	ctx, err := newCtx(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = load(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(ctx)
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
