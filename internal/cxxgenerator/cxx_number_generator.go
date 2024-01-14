package cxxgenerator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"text/template"
)

var cxxIntTypes = map[datatype.Kind]string{
	datatype.Int8:  "int8_t",
	datatype.Int32: "int32_t",
	datatype.Int64: "int64_t",
}

// CxxNumberGenerator is a DataTypeMessageCodeGenerator for NanoPack number types.
type CxxNumberGenerator struct{}

func (g CxxNumberGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return cxxIntTypes[dataType.Kind]
}

func (g CxxNumberGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g CxxNumberGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return cxxIntTypes[field.Type.Kind] + " " + strcase.ToSnake(field.Name)
}

func (g CxxNumberGenerator) ConstructorFieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g CxxNumberGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g CxxNumberGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		g.ReadValueFromBuffer(field.Type, s, ctx),
		fmt.Sprintf("this->%v = %v;", s, s))
}

func (g CxxNumberGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
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

func (g CxxNumberGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.write_field_size(%d, %d);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("writer.append_%v(%v);", field.Type.Identifier, strcase.ToSnake(field.Name)))
}

func (g CxxNumberGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("writer.append_%v(%v);", dataType.Identifier, varName)
}

func (g CxxNumberGenerator) ToFuncMap() template.FuncMap {
	return template.FuncMap{
		generator.FuncMapKeyTypeDeclaration:             g.TypeDeclaration,
		generator.FuncMapKeyReadSizeExpression:          g.ReadSizeExpression,
		generator.FuncMapKeyConstructorFieldParameter:   g.ConstructorFieldParameter,
		generator.FuncMapKeyConstructorFieldInitializer: g.ConstructorFieldInitializer,
		generator.FuncMapKeyFieldDeclaration:            g.FieldDeclaration,
		generator.FuncMapKeyReadFieldFromBuffer:         g.ReadFieldFromBuffer,
		generator.FuncMapKeyReadValueFromBuffer:         g.ReadValueFromBuffer,
		generator.FuncMapKeyWriteFieldToBuffer:          g.WriteFieldToBuffer,
		generator.FuncMapKeyWriteVariableToBuffer:       g.WriteVariableToBuffer,
	}
}
