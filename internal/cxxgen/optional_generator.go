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
	gm cxxCodeFragmentGeneratorMap
}

func (g optionalGenerator) TypeDeclaration(dataType datatype.DataType) string {
	ig := g.gm[dataType.ElemType.Kind]
	if dataType.ElemType.Kind == datatype.Message && dataType.Schema == nil {
		return "std::unique_ptr<NanoPack::Message>"
	}
	return fmt.Sprintf("std::optional<%v>", ig.TypeDeclaration(*dataType.ElemType))
}

func (g optionalGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	if dataType.ElemType.Kind == datatype.Message && dataType.Schema == nil {
		return fmt.Sprintf("std::unique_ptr<NanoPack::Message> %v", paramName)
	}
	declr := g.TypeDeclaration(dataType)
	// TODO: extremely hacky way to detect whether a data type needs unique ptr
	// in order to determine whether a const should be inserted, but ig it works for now
	if strings.Contains(declr, "std::unique_ptr") {
		return fmt.Sprintf("%v &%v", declr, paramName)
	}
	return fmt.Sprintf("const %v &%v", declr, paramName)
}

func (g optionalGenerator) RValue(dataType datatype.DataType, argName string) string {
	if dataType.ElemType.Kind == datatype.Message && dataType.Schema == nil {
		return fmt.Sprintf("std::move(%v)", argName)
	}
	return argName
}

func (g optionalGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.ElemType.Kind == datatype.Message && dataType.Schema == nil {
		return g.gm[dataType.ElemType.Kind].ReadSizeExpression(*dataType.ElemType, varName)
	}
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

	var assign string
	if field.Type.ElemType.Kind == datatype.Message && field.Type.ElemType.Schema == nil {
		assign = fmt.Sprintf("    %v = nullptr;", s)
	} else {
		assign = fmt.Sprintf("    %v = std::nullopt;", s)
	}

	return generator.Lines(
		fmt.Sprintf("if (reader.read_field_size(%d) < 0) {", field.Number),
		assign,
		"} else {",
		ig.ReadFieldFromBuffer(field, ctx),
		"}")
}

func (g optionalGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ctx.AddVariableToScope(varName)
	ig := g.gm[dataType.ElemType.Kind]

	var initValue string
	if dataType.ElemType.Kind == datatype.Message && dataType.ElemType.Schema == nil {
		initValue = "nullptr"
	} else {
		initValue = "std::nullopt"
	}

	return generator.Lines(
		fmt.Sprintf("%v %v = %v;", g.TypeDeclaration(dataType), varName, initValue),
		"if (reader.buffer[ptr++] != 0) {",
		ig.ReadValueFromBuffer(dataType, varName, ctx),
		"}")
}

func (g optionalGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	ig := g.gm[field.Type.ElemType.Kind]

	var ifExpr string
	if field.Type.ElemType.Kind == datatype.Message && field.Type.ElemType.Schema == nil {
		ifExpr = fmt.Sprintf("%v == nullptr", s)
	} else {
		ifExpr = fmt.Sprintf("%v.has_value()", s)
	}

	return generator.Lines(
		fmt.Sprintf("if (%v) {", ifExpr),
		ig.WriteFieldToBuffer(field, ctx),
		"} else {",
		fmt.Sprintf("writer.write_field_size(%d, -1, offset);", field.Number),
		"}")
}

func (g optionalGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]

	var ifExpr string
	var rvalue string
	if dataType.ElemType.Kind == datatype.Message && dataType.ElemType.Schema == nil {
		ifExpr = fmt.Sprintf("%v != nullptr", varName)
		rvalue = ig.RValue(*dataType.ElemType, varName)
	} else {
		ifExpr = fmt.Sprintf("%v.has_value()", varName)
		rvalue = ig.RValue(*dataType.ElemType, varName+".value()")
	}

	return generator.Lines(
		fmt.Sprintf("if (%v) {", ifExpr),
		"    writer.append_uint8(1);",
		fmt.Sprintf("    const auto %v_value = %v;", varName, rvalue),
		ig.WriteVariableToBuffer(dataType, varName+"_value", ctx),
		"} else {",
		"    writer.append_uint8(0);",
		"}")
}
