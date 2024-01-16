package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type boolGenerator struct{}

func (g boolGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "Bool"
}

func (g boolGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return "1"
}

func (g boolGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g boolGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g boolGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g boolGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("let %v = data.read(at: ptr)", strcase.ToCamel(field.Name)),
		"ptr += 1")
}

func (g boolGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l0 string
	if ctx.IsVariableInScope(varName) {
		l0 = fmt.Sprintf("%v = data.read(at: ptr)", varName)
	} else {
		l0 = fmt.Sprintf("let %v: %v = data.read(at: ptr)", varName, g.TypeDeclaration(dataType))
	}
	return generator.Lines(
		l0,
		"ptr += 1")
}

func (g boolGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("data.write(size: 1, ofField: %d)", field.Number),
		fmt.Sprintf("data.append(bool: %v)", strcase.ToCamel(field.Name)))
}

func (g boolGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("data.append(bool: %v)", varName)
}
