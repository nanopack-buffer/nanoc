package cxxgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

var cxxIntTypes = map[datatype.Kind]string{
	datatype.Int8:   "int8_t",
	datatype.Int32:  "int32_t",
	datatype.Int64:  "int64_t",
	datatype.Double: "double",
	datatype.UInt8:  "uint8_t",
	datatype.UInt32: "uint32_t",
	datatype.UInt64: "uint64_t",
}

// numberGenerator is a DataTypeMessageCodeGenerator for NanoPack number types.
type numberGenerator struct{}

func (g numberGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return cxxIntTypes[dataType.Kind]
}

func (g numberGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g numberGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return cxxIntTypes[field.Type.Kind] + " " + strcase.ToSnake(field.Name)
}

func (g numberGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g numberGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g numberGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.Type.Kind == datatype.Enum {
		return g.ReadValueFromBuffer(*field.Type.ElemType, s+"_raw_value", ctx)
	}

	return generator.Lines(
		g.ReadValueFromBuffer(field.Type, s, ctx),
		fmt.Sprintf("this->%v = %v;", s, s))
}

func (g numberGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l string
	if ctx.IsVariableInScope(varName) {
		l = fmt.Sprintf("%v = reader.read_%v(ptr);", varName, dataType.Identifier)
	} else {
		l = fmt.Sprintf("const %v %v = reader.read_%v(ptr);", g.TypeDeclaration(dataType), varName, dataType.Identifier)
	}
	return generator.Lines(
		l,
		fmt.Sprintf("ptr += %d;", dataType.ByteSize))
}

func (g numberGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	if field.Type.Kind == datatype.Enum {
		return generator.Lines(
			fmt.Sprintf("NanoPack::write_field_size(%d, %d, offset, buf);", field.Number, field.Type.ByteSize),
			fmt.Sprintf("NanoPack::append_%v(%v.value(), buf);", field.Type.ElemType.Identifier, strcase.ToSnake(field.Name)),
			fmt.Sprintf("bytes_written += %d;", field.Type.ElemType.ByteSize))
	}
	return generator.Lines(
		fmt.Sprintf("NanoPack::write_field_size(%d, %d, offset, buf);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("NanoPack::append_%v(%v, buf);", field.Type.Identifier, strcase.ToSnake(field.Name)),
		fmt.Sprintf("bytes_written += %d;", field.Type.ByteSize))
}

func (g numberGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if dataType.Kind == datatype.Enum {
		return fmt.Sprintf("NanoPack::append_%v(%v.value(), buf);", dataType.ElemType.Identifier, varName)
	}
	return fmt.Sprintf("NanoPack::append_%v(%v, buf);", dataType.Identifier, varName)
}
