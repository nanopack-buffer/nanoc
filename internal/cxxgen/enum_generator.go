package cxxgen

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
	if dataType.ElemType.Kind == datatype.String {
		return fmt.Sprintf("%v.value().size()", varName)
	}
	ig := g.gm[dataType.ElemType.Kind]
	return ig.ReadSizeExpression(*dataType.ElemType, varName)
}

func (g enumGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("const %v &%v", g.TypeDeclaration(field.Type), strcase.ToSnake(field.Name))
}

func (g enumGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g enumGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g enumGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		ig.ReadFieldFromBuffer(field, ctx),
		fmt.Sprintf("%v = %v(%v_raw_value);", s, g.TypeDeclaration(field.Type), s))
}

func (g enumGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	t := g.TypeDeclaration(dataType)

	var l1 string
	if ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("%v = %v(%v_raw_value);", varName, t, varName)
	} else {
		l1 = fmt.Sprintf("%v %v = %v(%v_raw_value);", t, varName, t, varName)
	}

	return generator.Lines(
		ig.ReadValueFromBuffer(*dataType.ElemType, varName+"_raw_value", ctx),
		l1)
}

func (g enumGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	return ig.WriteFieldToBuffer(field, ctx)
}

func (g enumGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.WriteVariableToBuffer(dataType, varName, ctx)
}
