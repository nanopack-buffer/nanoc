package npschema

import "nanoc/internal/datatype"

type Enum struct {
	Name      string
	ValueType datatype.DataType
	Members   []EnumMember
}

type PartialEnum struct {
	Name          string
	ValueTypeName string
	Members       []EnumMember
}

type EnumMember struct {
	Name         string
	ValueLiteral string
}

func (e Enum) isSchema()               {}
func (p PartialEnum) isPartialSchema() {}
