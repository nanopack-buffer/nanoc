package tsgen

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
	return fmt.Sprintf("%v | null", ig.TypeDeclaration(*dataType.ElemType))
}

func (g optionalGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("%v ? %v : 1", varName, ig.ReadSizeExpression(*dataType.ElemType, varName))
}

func (g optionalGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g optionalGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v", c, c)
}

func (g optionalGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g optionalGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]

	ctx.AddVariableToScope(c)

	ls := generator.Lines(
		fmt.Sprintf("let %v: %v;", c, g.TypeDeclaration(field.Type)),
		fmt.Sprintf("if (reader.readFieldSize(%d) >= 0) {", field.Number),
		ig.ReadFieldFromBuffer(npschema.MessageField{
			Name:   field.Name,
			Type:   *field.Type.ElemType,
			Number: field.Number,
			Schema: field.Schema,
		}, ctx),
		"} else {",
		fmt.Sprintf("    %v = null;", c),
		"}")

	ctx.RemoveVariableFromScope(c)

	return ls
}

func (g optionalGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.Kind]

	ctx.AddVariableToScope(varName)

	ls := generator.Lines(
		fmt.Sprintf("let %v: %v;", varName, g.TypeDeclaration(dataType)),
		"if (reader.readBoolean(ptr++)) {",
		ig.ReadValueFromBuffer(*dataType.ElemType, varName, ctx),
		"} else {",
		fmt.Sprintf("%v = null;", varName))

	ctx.RemoveVariableFromScope(varName)

	return ls
}

func (g optionalGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]

	return generator.Lines(
		fmt.Sprintf("if (this.%v) {", c),
		ig.WriteFieldToBuffer(npschema.MessageField{
			Name:   field.Name,
			Type:   *field.Type.ElemType,
			Number: field.Number,
			Schema: field.Schema,
		}, ctx),
		"} else {",
		fmt.Sprintf("writer.writeFieldSize(%d, -1, offset);", field.Number),
		"}")
}

func (g optionalGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.Kind]
	return generator.Lines(
		fmt.Sprintf("if (%v) {", varName),
		"    writer.appendBoolean(true);",
		ig.WriteVariableToBuffer(*dataType.ElemType, varName, ctx),
		"} else {",
		"    writer.appendBoolean(false);",
		"}")
}
