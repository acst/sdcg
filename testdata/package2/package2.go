package package2

// StructC ...
type StructC struct {
	Field1 string
	Field2 []string
	Field3 map[string]string
	Field4 bool
	Field5 []*StructD
	Field6 *StructD
	Field7 map[string]*StructD
	Field8 string
	Filed9 string
}

// StructD ...
type StructD struct {
	Field1 string
	Field2 bool
}

type PublicStructA struct {
	Field1 string
	Field2 bool
}
