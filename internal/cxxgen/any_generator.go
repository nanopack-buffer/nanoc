package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type anyGenerator struct{}

func (g anyGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "NanoPack::Any"
}

func (g anyGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return varName + ".size()"
}

func (g anyGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return "const NanoPack::Any &" + paramName
}

func (g anyGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g anyGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g anyGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g anyGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("%v %v;", g.TypeDeclaration(field.Type), strcase.ToSnake(field.Name))
}

func (g anyGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		fmt.Sprintf("const int32_t %v_byte_size = reader.read_field_size(%d);", s, field.Number),
		fmt.Sprintf("%v = %v(reader.buffer + ptr, %v_byte_size);", s, g.TypeDeclaration(field.Type), s),
		fmt.Sprintf("ptr += %v_byte_size;", s))
}

func (g anyGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l2 string
	if ctx.IsVariableInScope(varName) {
		l2 = fmt.Sprintf("%v = %v(reader.buffer + ptr, %v_byte_size);", varName, g.TypeDeclaration(dataType), varName)
	} else {
		l2 = fmt.Sprintf("%v %v(reader.buffer + ptr, %v_byte_size);", g.TypeDeclaration(dataType), varName, varName)
	}

	return generator.Lines(
		fmt.Sprintf("const int32_t %v_byte_size;", varName),
		fmt.Sprintf("reader.read_int32(ptr, %v_byte_size);", varName),
		"ptr += 4;",
		l2,
		fmt.Sprintf("ptr += %v_byte_size;", varName))
}

func (g anyGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	return generator.Lines(
		fmt.Sprintf("writer.write_field_size(%d, %v.size(), offset);", field.Number, s),
		fmt.Sprintf("writer.append_bytes(%v.data(), %v.size());", s, s))
}

func (g anyGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.append_int32(%v.size());", varName),
		fmt.Sprintf("writer::append_bytes(%v.data(), %v.size());", varName, varName))
}
