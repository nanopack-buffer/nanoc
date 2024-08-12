package npschema

import "nanoc/internal/datatype"

type Enum struct {
	IsDefaultValueUsed bool
	SchemaPath         string
	Name               string
	ValueType          datatype.DataType
	Members            []EnumMember
}

type PartialEnum struct {
	IsDefaultValueUsed bool
	SchemaPath         string
	Name               string
	ValueTypeName      string
	Members            []EnumMember
}

type EnumMember struct {
	Name         string
	ValueLiteral string
}

func (e *Enum) SchemaPathAbsolute() string {
	return e.SchemaPath
}

func (e *Enum) DataType() *datatype.DataType {
	return &datatype.DataType{
		Identifier: e.Name,
		Kind:       datatype.Enum,
		ByteSize:   e.ValueType.ByteSize,
		Schema:     e,
		KeyType:    nil,
		ElemType:   &e.ValueType,
	}
}

func (e *Enum) SchemaName() string {
	return e.Name
}

func (p PartialEnum) IsPartialSchema() {}
