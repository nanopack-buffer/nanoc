package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type boolGenerator struct{}

func (g boolGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "bool"
}

func (g boolGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return "bool " + paramName
}

func (g boolGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g boolGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g boolGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g boolGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g boolGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g boolGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return fmt.Sprintf("reader.read_bool(ptr++, %v);", strcase.ToSnake(field.Name))
}

func (g boolGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var declr string
	if !ctx.IsVariableInScope(varName) {
		declr = fmt.Sprintf("%v %v;", g.TypeDeclaration(dataType), varName)
	}
	return generator.Lines(
		declr,
		fmt.Sprintf("reader.read_bool(ptr++, %v);", varName))
}

func (g boolGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.write_field_size(%d, %d, offset);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("writer.append_bool(%v);", strcase.ToSnake(field.Name)))
}

func (g boolGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("writer.append_bool(%v);", varName)
}
