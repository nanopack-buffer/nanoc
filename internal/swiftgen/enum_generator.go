package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type enumGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g enumGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return dataType.Identifier
}

func (g enumGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.ReadSizeExpression(dataType, varName)
}

func (g enumGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g enumGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g enumGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g enumGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	c := strcase.ToLowerCamel(field.Name)

	var tmpv string
	var l4 string
	if ctx.IsVariableInScope(c) {
		tmpv = c + "_"
		l4 = fmt.Sprintf("%v = %v", c, tmpv)
	} else {
		tmpv = c
	}

	return generator.Lines(
		ig.ReadFieldFromBuffer(field, ctx),
		fmt.Sprintf("guard let %v = %v(rawValue: %vRawValue) else {", tmpv, g.TypeDeclaration(field.Type), c),
		"    return nil",
		"}",
		l4)
}

func (g enumGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]

	var tmpv string
	var l4 string
	if ctx.IsVariableInScope(varName) {
		tmpv = varName + "_"
		l4 = fmt.Sprintf("%v = %v", varName, tmpv)
	} else {
		tmpv = varName
	}

	return generator.Lines(
		ig.ReadValueFromBuffer(dataType, varName, ctx),
		fmt.Sprintf("guard let %v = %v(rawValue: %vRawValue) else {", tmpv, g.TypeDeclaration(dataType), varName),
		"    return nil",
		"}",
		l4)
}

func (g enumGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	return ig.WriteFieldToBuffer(field, ctx)
}

func (g enumGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.WriteVariableToBuffer(dataType, varName, ctx)
}
