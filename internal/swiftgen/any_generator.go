package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type anyGenerator struct{}

func (g anyGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "Data"
}

func (g anyGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v.count", varName)
}

func (g anyGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g anyGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g anyGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g anyGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)

	var l1 string
	if ctx.IsVariableInScope(c) {
		l1 = fmt.Sprintf("%v = data[ptr..<ptr + %vByteSize]", c, c)
	} else {
		l1 = fmt.Sprintf("let %v = data[ptr..<ptr + %vByteSize]", c, c)
	}

	return generator.Lines(
		fmt.Sprintf("let %vByteSize = data.readSize(ofField: %d)", c, field.Number),
		l1,
		fmt.Sprintf("ptr += %vByteSize", c))
}

func (g anyGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l2 string
	if ctx.IsVariableInScope(varName) {
		l2 = fmt.Sprintf("%v = data[ptr..<ptr + %vByteSize]", varName, varName)
	} else {
		l2 = fmt.Sprintf("let %v = data[ptr..<ptr + %vByteSize]", varName, varName)
	}

	return generator.Lines(
		fmt.Sprintf("let %vByteSize = data.readSize(at: ptr)"),
		"ptr += 4",
		l2,
		fmt.Sprintf("ptr += %vByteSize", varName))
}

func (g anyGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("data.write(size: %v.count, ofField: %d)", c, field.Number),
		fmt.Sprintf("data.append(%v)", c))
}

func (g anyGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("data.append(size: %v.count)", varName),
		fmt.Sprintf("data.append(%v)", varName))
}
