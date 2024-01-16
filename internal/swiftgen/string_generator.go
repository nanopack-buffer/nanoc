package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type stringGenerator struct{}

func (g stringGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "String"
}

func (g stringGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v.lengthOfBytes(using .utf8)", varName)
}

func (g stringGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g stringGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g stringGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g stringGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("let %vSize = data.readSize(ofField: %d)", c, field.Number),
		fmt.Sprintf("guard let %v = data.read(at: ptr, withLength: %vSize) else {", c, c),
		"    return nil",
		"}",
		fmt.Sprintf("ptr += %vSize", c))
}

func (g stringGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("let %vSize = readSize(at: ptr)", varName),
		"    ptr += 4",
		fmt.Sprintf("guard let %v = data.read(at: ptr, withLength: %vSize) else {", varName, varName),
		"    return nil",
		"}",
		fmt.Sprintf("ptr += %vSize", varName))
}

func (g stringGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("data.write(size: %v, ofField: %d)", g.ReadSizeExpression(field.Type, c), field.Number),
		fmt.Sprintf("data.append(string: %v)", c))
}

func (g stringGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("data.append(size: %v)", g.ReadSizeExpression(dataType, varName)),
		fmt.Sprintf("data.append(string: %v)", varName))
}
