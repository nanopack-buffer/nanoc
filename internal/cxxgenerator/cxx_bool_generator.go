package cxxgenerator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"text/template"
)

type cxxBoolGenerator struct{}

func (g cxxBoolGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "bool"
}

func (g cxxBoolGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g cxxBoolGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g cxxBoolGenerator) ConstructorFieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g cxxBoolGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g cxxBoolGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		g.ReadValueFromBuffer(field.Type, s, ctx),
		fmt.Sprintf("this->%v = %v", s, s))
}

func (g cxxBoolGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if ctx.IsVariableInScope(varName) {
		return fmt.Sprintf("%v = reader.read_bool(ptr++);", varName)
	}
	return fmt.Sprintf("const %v %v = reader.read_bool(ptr++);", g.TypeDeclaration(dataType), varName)
}

func (g cxxBoolGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.write_field_size(%d, %d);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("writer.append_bool(%v);", strcase.ToSnake(field.Name)))
}

func (g cxxBoolGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("writer.append_bool(%v);", varName)
}

func (g cxxBoolGenerator) ToFuncMap() template.FuncMap {
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
