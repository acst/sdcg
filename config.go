package main

type Package struct {
	Name     string `hcl:"name,label"`
	FileName string `hcl:"filename,optional"`
}

type Use struct {
	Name  string `hcl:"name,label"`
	Alias string `hcl:"alias,optional"`
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
