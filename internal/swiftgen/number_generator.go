package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"text/template"
)

type numberGenerator struct{}

var intTypes = map[datatype.Kind]string{
	datatype.Int8:   "Int8",
	datatype.Int32:  "Int32",
	datatype.Int64:  "Int64",
	datatype.Double: "Double",
}

func (g numberGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return intTypes[dataType.Kind]
}

func (g numberGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g numberGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g numberGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g numberGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g numberGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return g.ReadValueFromBuffer(field.Type, strcase.ToCamel(field.Name), ctx)
}

func (g numberGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l0 string
	if ctx.IsVariableInScope(varName) {
		l0 = fmt.Sprintf("%v = data.read(at: ptr)", varName)
	} else {
		l0 = fmt.Sprintf("let %v: %v = data.read(at: ptr)", varName, g.TypeDeclaration(dataType))
	}
	return generator.Lines(
		l0,
		fmt.Sprintf("ptr += %d", dataType.ByteSize))
}

func (g numberGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	var l1 string
	if field.Type.Kind == datatype.Double {
		l1 = fmt.Sprintf("data.append(double: %v)", strcase.ToSnake(field.Name))
	} else {
		l1 = fmt.Sprintf("data.append(int: %v", strcase.ToSnake(field.Name))
	}
	return generator.Lines(
		fmt.Sprintf("data.write(size: %d, ofField: %d)", field.Type.ByteSize, field.Number),
		l1)
}

func (g numberGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if dataType.Kind == datatype.Double {
		return fmt.Sprintf("data.append(double: %v)", varName)
	}
	return fmt.Sprintf("data.append(int: %v)", varName)
}

func (g numberGenerator) ToFuncMap() template.FuncMap {
	//TODO implement me
	panic("implement me")
}
