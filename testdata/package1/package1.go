package package1

// Enum ...
type Enum int

const (
	EnumVal1 Enum = iota
	EnumVal2
	EnumVal3
)

// StructA ...
type StructA struct {
	Field1 string
	Field2 bool
	Field3 map[string]string
	Field4 []string
	Field5 []*StructB
	Field6 *StructB
	Field7 map[string]*StructB
	Field8 Enum
	Field9 []string
}

// StructB ...
type StructB struct {
	Field1 bool
	Field2 string
}

type privateStructA struct {
	field1 string
	field2 bool
}

type privateStructB struct {
	field1 string
	field2 bool
}
