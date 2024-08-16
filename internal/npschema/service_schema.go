package npschema

import "nanoc/internal/datatype"

type PartialService struct {
	SchemaPath        string
	Name              string
	ImportedTypeNames []string
	DeclaredFunctions []PartialDeclaredFunction
}

type Service struct {
	SchemaPath        string
	Name              string
	ImportedTypes     []datatype.DataType
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

	// ParametersByteSize stores the total number of bytes required to store all of parameters.
	// This is datatype.DynamicSize if any of the parameter does not have a static size.
	ParametersByteSize int
	ReturnType         *datatype.DataType
	ErrorType          *datatype.DataType
}

type PartialFunctionParameter struct {
	Name     string
	TypeName string
}

type FunctionParameter struct {
	Name string
	Type datatype.DataType
}

func (s PartialService) IsPartialSchema() {}

func (s *Service) SchemaPathAbsolute() string {
	return s.SchemaPath
}

func (s *Service) DataType() *datatype.DataType {
	return nil
}

func (s *Service) SchemaName() string {
	return s.Name
}
