package tsgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type boolGenerator struct{}

func (g boolGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "boolean"
}

func (g boolGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return "1"
}

func (g boolGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g boolGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g boolGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g boolGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return g.ReadValueFromBuffer(field.Type, strcase.ToLowerCamel(field.Name), ctx)
}

func (g boolGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if ctx.IsVariableInScope(varName) {
		return fmt.Sprintf("%v = reader.readBoolean(ptr++);", varName)
	}
	return fmt.Sprintf("const %v = reader.readBoolean(ptr++);", varName)
}

func (g boolGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.appendBoolean(this.%v);", strcase.ToLowerCamel(field.Name)),
		fmt.Sprintf("writer.writeFieldSize(%d, 1, offset);", field.Number),
		"bytesWritten += 1;")
}

func (g boolGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("writer.appendBoolean(%v);", varName)
}
