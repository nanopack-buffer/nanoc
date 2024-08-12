package npschema

import "nanoc/internal/datatype"

type PartialServiceSchema struct {
	SchemaPath        string
	Name              string
	ImportedTypeNames []string
	DeclaredFunctions []PartialDeclaredFunction
}

type ServiceSchema struct {
	SchemaPath        string
	Name              string
	ImportedTypes     []datatype.Schema
	DeclaredFunctions []DeclaredFunction
}

type PartialDeclaredFunction struct {
	Name           string
	Parameters     []PartialFunctionParameter
	ReturnTypeName string
	ErrorTypeName  string
}

type DeclaredFunction struct {
	Name       string
	Parameters []FunctionParameter
	ReturnType *datatype.DataType
	ErrorType  *datatype.DataType
}

type PartialFunctionParameter struct {
	Name     string
	TypeName string
}

type FunctionParameter struct {
	Name string
	Type datatype.DataType
}

func (s PartialServiceSchema) IsPartialSchema() {}

func (s *ServiceSchema) SchemaPathAbsolute() string {
	return s.SchemaPath
}

func (s *ServiceSchema) DataType() *datatype.DataType {
	return nil
}

func (s *ServiceSchema) SchemaName() string {
	return s.Name
}
