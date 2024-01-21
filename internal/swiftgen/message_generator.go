package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type messageGenerator struct{}

func (g messageGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return dataType.Identifier
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%vData.count", varName)
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	var v string
	var l4 string
	if ctx.IsVariableInScope(c) {
		v = c + "_"
		l4 = fmt.Sprintf("%v = %v", c, v)
	} else {
		v = c
	}

	return generator.Lines(
		fmt.Sprintf("let %vByteSize = data.readSize(ofField: %d)", c, field.Number),
		fmt.Sprintf("guard let %v = %v(data: data[ptr...]) else {", v, g.TypeDeclaration(field.Type)),
		"    return nil",
		"}",
		l4,
		fmt.Sprintf("ptr += %vByteSize", c))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var v string
	var l4 string
	if ctx.IsVariableInScope(varName) {
		v = varName + "_"
		l4 = fmt.Sprintf("%v = %v", varName, v)
	} else {
		v = varName
	}

	return generator.Lines(
		fmt.Sprintf("var %vByteSize = 0", varName),
		fmt.Sprintf("guard let %v = %v(data: data[ptr...], bytesRead: &%vByteSize) else {", v, g.TypeDeclaration(dataType), varName),
		"    return nil",
		"}",
		l4,
		fmt.Sprintf("ptr += %vByteSize", varName))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("guard let %vData = %v.data() else {", c, c),
		"    return nil",
		"}",
		fmt.Sprintf("data.write(size: %vData.count, ofField: %d)", c, field.Number),
		fmt.Sprintf("data.append(%vData)", c))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("guard let %vData = %v.data() else {", varName, varName),
		"    return nil",
		"}",
		fmt.Sprintf("data.append(%vData)", varName))
}
