package tsgen

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
	return "T" + dataType.Identifier
}

func (g enumGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.ReadSizeExpression(*dataType.ElemType, varName)
}

func (g enumGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g enumGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g enumGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g enumGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	return ig.ReadFieldFromBuffer(field, ctx)
}

func (g enumGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.ReadValueFromBuffer(dataType, varName, ctx)
}

func (g enumGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	return ig.WriteFieldToBuffer(npschema.MessageField{
		Name:   field.Name,
		Type:   *field.Type.ElemType,
		Number: field.Number,
		Schema: field.Schema,
	}, ctx)
}

func (g enumGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	return ig.WriteVariableToBuffer(*dataType.ElemType, varName, ctx)
}
