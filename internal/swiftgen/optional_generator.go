package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type optionalGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g optionalGenerator) TypeDeclaration(dataType datatype.DataType) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.TypeDeclaration(*dataType.ElemType) + "?"
}

func (g optionalGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("%v.map { %v } ?? 1", varName, ig.ReadSizeExpression(*dataType.ElemType, "$0"))
}

func (g optionalGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g optionalGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g optionalGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g optionalGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]

	var l0 string
	if !ctx.IsVariableInScope(c) {
		l0 = fmt.Sprintf("var %v: %v", c, g.TypeDeclaration(field.Type))
	}
	ctx.AddVariableToScope(c)

	return generator.Lines(
		l0,
		fmt.Sprintf("if data.readSize(ofField: %d) < 0 {", field.Number),
		fmt.Sprintf("    %v = nil", c),
		"} else {",
		ig.ReadFieldFromBuffer(npschema.MessageField{
			Name:   field.Name,
			Type:   *field.Type.ElemType,
			Number: field.Number,
			Schema: field.Schema,
		}, ctx),
		"}")
}

func (g optionalGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]

	var l0 string
	if !ctx.IsVariableInScope(varName) {
		l0 = fmt.Sprintf("var %v: %v", varName, g.TypeDeclaration(dataType))
	}
	ctx.AddVariableToScope(varName)

	return generator.Lines(
		l0,
		"if data[ptr] != 0 {",
		"    ptr += 1",
		ig.ReadValueFromBuffer(*dataType.ElemType, varName, ctx),
		"} else {",
		"    ptr += 1",
		fmt.Sprintf("%v = nil", varName),
		"}")
}

func (g optionalGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]
	return generator.Lines(
		fmt.Sprintf("if let %v = self.%v {", c, c),
		ig.WriteFieldToBuffer(npschema.MessageField{
			Name:   field.Name,
			Type:   *field.Type.ElemType,
			Number: field.Number,
			Schema: field.Schema,
		}, ctx),
		"} else {",
		fmt.Sprintf("data.write(size: -1, ofField: %d)", field.Number),
		"}")
}

func (g optionalGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return generator.Lines(
		fmt.Sprintf("if let %v {", varName),
		"    data.append(byte: 1)",
		ig.WriteVariableToBuffer(*dataType.ElemType, varName, ctx),
		"} else {",
		"    data.append(byte: 0)",
		"}")
}
