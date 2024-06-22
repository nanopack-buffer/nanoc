package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strings"

	"github.com/iancoleman/strcase"
)

type optionalGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g optionalGenerator) TypeDeclaration(dataType datatype.DataType) string {
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("std::optional<%v>", ig.TypeDeclaration(*dataType.ElemType))
}

func (g optionalGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	ig := g.gm[dataType.ElemType.Kind]
	expr := ig.ReadSizeExpression(*dataType.ElemType, varName)
	// if the expression uses dot to access varName, e.g. varName.size(), replace dot with ->
	expr = strings.Replace(expr, varName+".", varName+"->", -1)
	return fmt.Sprintf("%v.has_value() ? %v : 1", varName, expr)
}

func (g optionalGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g optionalGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g optionalGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g optionalGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		fmt.Sprintf("if (reader.read_field_size(%d) < 0) {", field.Number),
		fmt.Sprintf("    %v = std::nullopt;", s),
		"} else {",
		ig.ReadFieldFromBuffer(field, ctx),
		"}")
}

func (g optionalGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ctx.AddVariableToScope(varName)
	ig := g.gm[dataType.ElemType.Kind]
	return generator.Lines(
		fmt.Sprintf("%v %v = std::nullopt;", g.TypeDeclaration(dataType), varName),
		"if (reader.buffer[ptr++] != 0) {",
		ig.ReadValueFromBuffer(dataType, varName, ctx),
		"}")
}

func (g optionalGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]
	return generator.Lines(
		fmt.Sprintf("if (%v.has_value()) {", s),
		ig.WriteFieldToBuffer(field, ctx),
		"} else {",
		fmt.Sprintf("writer.write_field_size(%d, -1, offset);", field.Number),
		"}")
}

func (g optionalGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return generator.Lines(
		fmt.Sprintf("if (%v.has_value()) {", varName),
		"    writer.append_uint8(1);",
		fmt.Sprintf("    const auto %v_value = %v.value();", varName, varName),
		ig.WriteVariableToBuffer(*dataType.ElemType, varName+"_value", ctx),
		"} else {",
		"    writer.append_uint8(0);",
		"}")
}
