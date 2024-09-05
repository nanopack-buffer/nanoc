package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type stringGenerator struct{}

func (g stringGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "std::string"
}

func (g stringGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return "const std::string &" + paramName
}

func (g stringGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g stringGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return varName + ".size() + 4"
}

func (g stringGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g stringGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g stringGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g stringGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	switch field.Type.Kind {
	case datatype.Enum:
		return generator.Lines(
			fmt.Sprintf("const int32_t %v_size = reader.read_field_size(%d);", s, field.Number),
			fmt.Sprintf("%v %v_raw_value;", g.TypeDeclaration(*field.Type.ElemType), s),
			fmt.Sprintf("reader.read_string(ptr, %v_size, %v_raw_value);", s, s),
			fmt.Sprintf("ptr += %v_size;", s))

	case datatype.Optional:
		return generator.Lines(
			fmt.Sprintf("%v = std::move(std::string());", s),
			fmt.Sprintf("const int32_t %v_size = reader.read_field_size(%d);", s, field.Number),
			fmt.Sprintf("reader.read_string(ptr, %v_size, %v);", s, s),
			fmt.Sprintf("ptr += %v_size;", s))

	default:
		return generator.Lines(
			fmt.Sprintf("const int32_t %v_size = reader.read_field_size(%d);", s, field.Number),
			fmt.Sprintf("reader.read_string(ptr, %v_size, %v);", s, s),
			fmt.Sprintf("ptr += %v_size;", s))
	}
}

func (g stringGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var declr string
	if !ctx.IsVariableInScope(varName) {
		declr = fmt.Sprintf("%v %v;", g.TypeDeclaration(dataType), varName)
	}

	return generator.Lines(
		fmt.Sprintf("uint32_t %v_size;", varName),
		fmt.Sprintf("reader.read_uint32(ptr, %v_size);", varName),
		"ptr += 4;",
		declr,
		fmt.Sprintf("reader.read_string(ptr, %v_size, %v);", varName, varName),
		fmt.Sprintf("ptr += %v_size;", varName))
}

func (g stringGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	switch field.Type.Kind {
	case datatype.Enum:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %v.value().size(), offset);", field.Number, s),
			fmt.Sprintf("writer.append_string(%v.value());", s))

	case datatype.Optional:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %v->size(), offset);", field.Number, s),
			fmt.Sprintf("writer.append_string(*%v);", s))

	default:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %v.size(), offset);", field.Number, s),
			fmt.Sprintf("writer.append_string(%v);", s))
	}
}

func (g stringGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var expr string
	if dataType.Kind == datatype.Enum {
		expr = varName + ".value()"
	} else {
		expr = varName
	}

	return generator.Lines(
		fmt.Sprintf("writer.append_int32(%v.size());", expr),
		fmt.Sprintf("writer.append_string(%v);", expr))
}
