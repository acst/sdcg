package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sirkon/goproxy/gomod"
	"golang.org/x/tools/go/packages"
)

type Package struct {
	Import   string `hcl:"import,label"`
	FileName string `hcl:"filename,optional"`
}

type Use struct {
	Import string `hcl:"import,label"`
	Alias  string `hcl:"alias,optional"`
}

type Convert struct {
	Type  string `hcl:"type,attr"`
	To    string `hcl:"to,attr"`
	Using string `hcl:"using,optional"`
}

type Field struct {
	Name  string `hcl:"name,label"`
	To    string `hcl:"to,optional"`
	Using string `hcl:"using,optional"`
}

type Struct struct {
	Name   string   `hcl:"name,label"`
	Ignore []string `hcl:"ignore,optional"`
}

type Map struct {
	From   *Struct  `hcl:"from,block"`
	To     *Struct  `hcl:"to,block"`
	Fields []*Field `hcl:"field,block"`
}

type File struct {
	Package  *Package   `hcl:"package,block"`
	Uses     []*Use     `hcl:"use,block"`
	Maps     []*Map     `hcl:"map,block"`
	Converts []*Convert `hcl:"convert,block"`
}

type sdcgContext struct {
	dir           string
	currentModule string

	// should not modify
	files []*contextFile
}

// If dir is not set, it defaults to the current working directory
func newCtx(dir string) (*sdcgContext, error) {
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("unable to get working directory: %w", err)
		}
	}
	modFile, err := findGoModFile(dir)
	if err != nil {
		return nil, err
	}
	mod, err := gomod.Parse("go.mod", modFile)
	if err != nil {
		return nil, err
	}
	return &sdcgContext{
		dir:           dir,
		currentModule: mod.Name,
	}, nil
}

func (ctx *sdcgContext) hasErrors() bool {
	for _, file := range ctx.files {
		if file.hasErrors() {
			return true
		}
	}
	return false
}

type contextFile struct {
	*File
	info   os.FileInfo
	errs   []error
	genPkg *packages.Package
	pkgs   map[string]*packages.Package
}

func (cf *contextFile) addErr(err error) {
	cf.errs = append(cf.errs, err)
}

func (cf *contextFile) hasErrors() bool {
	return len(cf.errs) != 0
}

func load(ctx *sdcgContext) error {
	err := parseHCLFiles(ctx)
	if err != nil {
		return err
	}
	loadPackages(ctx)
	// TODO: load converts
	// TODO: load maps
	return nil
}

func parseHCLFiles(ctx *sdcgContext) error {
	files, err := ioutil.ReadDir(ctx.dir)
	if err != nil {
		return fmt.Errorf("unable to read directory: %w", err)
	}
	parser := hclparse.NewParser()
	for _, f := range files {
		ctxFile := &contextFile{
			info: f,
		}
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".hcl" {
			continue
		}
		hcl, diag := parser.ParseHCLFile(path.Join(dir, f.Name()))
		if diag.HasErrors() {
			ctxFile.addErr(fmt.Errorf("unable to parse: %w", diag))
			ctx.files = append(ctx.files, ctxFile)
			continue
		}
		var config File
		diag = gohcl.DecodeBody(hcl.Body, hclCtx, &config)
		if diag.HasErrors() {
			ctxFile.addErr(fmt.Errorf("unable to decode into go struct: %w", diag))
			ctx.files = append(ctx.files, ctxFile)
			continue
		}
		ctxFile.File = &config
		ctx.files = append(ctx.files, ctxFile)
	}
	return nil
}

func loadPackages(ctx *sdcgContext) {
	for _, file := range ctx.files {
		if file.hasErrors() {
			continue
		}
		file.loadPackages(ctx)
	}
}

func (cf *contextFile) loadPackages(ctx *sdcgContext) {
	if cf.hasErrors() {
		cf.addErr(fmt.Errorf("cannot load packages from file with errors"))
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
		Dir:  dir,
	}

	// Load the package statement
	if cf.Package == nil {
		cf.addErr(fmt.Errorf("invalid file: must have package block"))
		return
	}
	if !strings.HasPrefix(cf.Package.Import, ctx.currentModule) {
		cf.addErr(fmt.Errorf("can only have a package statement apart of current module"))
		return
	}
	pkgImport, err := packages.Load(cfg, cf.Package.Import)
	if err != nil {
		cf.addErr(err)
		return
	}
	if len(pkgImport) != 1 {
		cf.addErr(fmt.Errorf("unexpected number of imported packages for package block"))
		return
	}
	if len(pkgImport[0].Errors) > 0 {
		for _, err := range pkgImport[0].Errors {
			if strings.Contains(err.Error(), "no matching versions for query") {
				// Package doesn't exist so we can create it
				continue
			}
			cf.addErr(fmt.Errorf("issues loading %s: %w", cf.Package.Import, err))
		}
		if cf.hasErrors() {
			return
		}
	} else {
		cf.genPkg = pkgImport[0]
	}

	// Load the Uses statements
	var pkgImports []string
	for _, u := range cf.Uses {
		pkgImports = append(pkgImports, u.Import)
	}
	loadedPkgs, err := packages.Load(cfg, pkgImports...)
	if err != nil {
		cf.addErr(err)
		return
	}
	for i, lp := range loadedPkgs {
		u := cf.Uses[i]
		if len(lp.Errors) != 0 {
			for _, err = range lp.Errors {
				cf.addErr(fmt.Errorf("issues loading %s: %w", u.Import, err))
			}
			continue
		}
		imp := lp.Name
		if u.Alias != "" {
			imp = u.Alias
		}
		if cf.pkgs == nil {
			cf.pkgs = make(map[string]*packages.Package)
		}
		cf.pkgs[imp] = lp
	}

	return
}

func findGoModFile(startingDir string) ([]byte, error) {
	dir, err := filepath.Abs(startingDir)
	if err != nil {
		return nil, err
	}
	for {
		if dir == string(filepath.Separator) {
			break
		}
		f, err := ioutil.ReadFile(path.Join(dir, "go.mod"))
		if errors.Is(err, os.ErrNotExist) {
			dir = path.Dir(dir)
			continue
		}
		return f, nil
	}
	return nil, fmt.Errorf("could not find go.mod file")
}
