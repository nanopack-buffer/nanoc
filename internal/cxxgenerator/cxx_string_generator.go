package cxxgenerator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"text/template"
)

type stringGenerator struct{}

func (g stringGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "std::string"
}

func (g stringGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return varName + ".size() + 4"
}

func (g stringGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g stringGenerator) ConstructorFieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g stringGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g stringGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		fmt.Sprintf("const int32_t %v_size = reader.read_field_size(%d);", s, field.Number),
		fmt.Sprintf("%v = reader.read_string(ptr, %v_size);", s, s),
		fmt.Sprintf("ptr += %v_size;", s))
}

func (g stringGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l2 string
	if ctx.IsVariableInScope(varName) {
		l2 = fmt.Sprintf("%v = reader.read_string(ptr, %v_size);", varName, varName)
	} else {
		l2 = fmt.Sprintf("%v %v = reader.read_string(ptr, %v_size);", g.TypeDeclaration(dataType), varName)
	}

	return generator.Lines(
		fmt.Sprintf("const int32_t %v_size = reader.read_int32(ptr);", varName),
		"ptr += 4;",
		l2,
		fmt.Sprintf("ptr += %v_size;", varName))
}

func (g stringGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		fmt.Sprintf("writer.write_field_size(%d, %v.size());", field.Number, s),
		fmt.Sprintf("writer.append_string(%v);", s))
}

func (g stringGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.append_int32(%v.size());", varName),
		fmt.Sprintf("writer.append_string(%v);", varName))
}

func (g stringGenerator) ToFuncMap() template.FuncMap {
	//TODO implement me
	panic("implement me")
}
