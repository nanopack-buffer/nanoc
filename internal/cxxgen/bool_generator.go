package cxxgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type boolGenerator struct{}

func (g boolGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "bool"
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
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		g.ReadValueFromBuffer(field.Type, s, ctx),
		fmt.Sprintf("this->%v = %v", s, s))
}

func (g boolGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if ctx.IsVariableInScope(varName) {
		return fmt.Sprintf("%v = reader.read_bool(ptr++);", varName)
	}
	return fmt.Sprintf("const %v %v = reader.read_bool(ptr++);", g.TypeDeclaration(dataType), varName)
}

func (g boolGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("NanoPack::write_field_size(%d, %d, offset, buf);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("NanoPack::append_bool(%v, buf);", strcase.ToSnake(field.Name)))
}

func (g boolGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("NanoPack::append_bool(%v, buf);", varName)
}
