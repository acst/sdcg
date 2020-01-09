package main

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

type goType uint64

const (
	gtInt goType = iota + 1<<32
	gtString
	gtFloat32
	gtFloat64
)

var stringGoTypeCty = map[string]cty.Value{
	"int":     cty.NumberUIntVal(uint64(gtInt)),
	"string":  cty.NumberUIntVal(uint64(gtString)),
	"float32": cty.NumberUIntVal(uint64(gtFloat32)),
	"float64": cty.NumberUIntVal(uint64(gtFloat64)),
}

type converterType uint64

const (
	ctStringToSlice converterType = iota + 1<<33
	ctSliceToString
)

var stringConverterTypeCty = map[string]cty.Value{
	"stringToSlice": cty.NumberUIntVal(uint64(ctStringToSlice)),
	"sliceToString": cty.NumberUIntVal(uint64(ctSliceToString)),
}

func joinCtyVariables(maps ...map[string]cty.Value) map[string]cty.Value {
	var totalSize int
	for _, m := range maps {
		totalSize += len(m)
	}
	final := make(map[string]cty.Value, totalSize)
	for _, m := range maps {
		for k, v := range m {
			final[k] = v
		}
	}
	return final
}

var stringFuncTypeCty = map[string]function.Function{
	"converter": function.New(&function.Spec{
		Params: []function.Parameter{
			function.Parameter{
				Name: "type",
				Type: cty.Number,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: converter,
	}),
	"gotype": function.New(&function.Spec{
		Params: []function.Parameter{
			function.Parameter{
				Name: "type",
				Type: cty.Number,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: gotype,
	}),
}

func converter(args []cty.Value, retType cty.Type) (cty.Value, error) {
	var val uint64
	err := gocty.FromCtyValue(args[0], &val)
	if err != nil {
		return cty.NullVal(cty.Number), fmt.Errorf("converter func: %w", err)
	}
	if val&uint64(1<<33) != uint64(1<<33) {
		return cty.NullVal(cty.Number), fmt.Errorf("converter was passed in unkown type")
	}
	return cty.NumberUIntVal(val), nil
}

func gotype(args []cty.Value, retType cty.Type) (cty.Value, error) {
	var val uint64
	err := gocty.FromCtyValue(args[0], &val)
	if err != nil {
		return cty.NullVal(cty.Number), fmt.Errorf("gotype func: %w", err)
	}
	if val&uint64(1<<32) != uint64(1<<32) {
		return cty.NullVal(cty.Number), fmt.Errorf("gotype was passed in unkown type")
	}
	return cty.NumberUIntVal(val), nil
}
