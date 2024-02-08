package cxxgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type messageGenerator struct{}

func (g messageGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return dataType.Identifier
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v_data.size()", varName)
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	ms := field.Type.Schema.(*npschema.Message)
	selfRef := ms.TypeID == field.Schema.TypeID
	td := g.TypeDeclaration(field.Type)
	s := strcase.ToSnake(field.Name)
	if selfRef {
		return fmt.Sprintf("std::shared_ptr<%v> %v", td, s)
	}
	return td + " " + s
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	ms := field.Type.Schema.(*npschema.Message)
	selfRef := ms.TypeID == field.Schema.TypeID
	s := strcase.ToSnake(field.Name)
	if selfRef {
		return fmt.Sprintf("%v(std::move(%v))", s, s)
	}
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	ms := field.Type.Schema.(*npschema.Message)
	selfRef := ms.TypeID == field.Schema.TypeID
	td := g.TypeDeclaration(field.Type)
	s := strcase.ToSnake(field.Name)
	if selfRef {
		return fmt.Sprintf("std::shared_ptr<%v> %v;", td, s)
	}
	return td + " " + s + ";"
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ms := field.Type.Schema.(*npschema.Message)
	selfRef := ms.TypeID == field.Schema.TypeID
	td := g.TypeDeclaration(field.Type)
	s := strcase.ToSnake(field.Name)

	if selfRef {
		return generator.Lines(
			fmt.Sprintf("if (reader.read_field_size(%d) < 0) {", field.Number),
			fmt.Sprintf("    %v = nullptr;", s),
			"} else {",
			fmt.Sprintf("    int %v_bytes_read = 0;", s),
			fmt.Sprintf("    %v = std::make_shared<%v>(begin + ptr, %v_bytes_read);", s, td, s),
			fmt.Sprintf("    ptr += %v_bytes_read;", s),
			"}")
	}

	return generator.Lines(
		fmt.Sprintf("int %v_bytes_read = 0;", s),
		fmt.Sprintf("%v = %v(begin + ptr, %v_bytes_read);", s, td, s),
		fmt.Sprintf("ptr += %v_bytes_read;", s))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	td := g.TypeDeclaration(dataType)
	var l1 string
	if ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("%v = %v(begin + ptr, %v_bytes_read);", varName, td, varName)
	} else {
		l1 = fmt.Sprintf("%v %v(begin + ptr, %v_bytes_read);", td, varName, varName)
	}
	return generator.Lines(
		fmt.Sprintf("int %v_bytes_read = 0;", varName),
		l1,
		fmt.Sprintf("ptr += %v_bytes_read;", varName))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ms := field.Type.Schema.(*npschema.Message)
	selfRef := ms.TypeID == field.Schema.TypeID
	s := strcase.ToSnake(field.Name)

	if selfRef {
		return generator.Lines(
			fmt.Sprintf("if (%v != nullptr) {", s),
			fmt.Sprintf("    const std::vector<uint8_t> %v_data = %v->data();", s, s),
			fmt.Sprintf("    writer.append_bytes(%v_data);", s),
			fmt.Sprintf("    writer.write_field_size(%d, %v_data.size());", field.Number, s),
			"} else {",
			fmt.Sprintf("    writer.write_field_size(%d, -1);", field.Number),
			"}")
	}

	return generator.Lines(
		fmt.Sprintf("const std::vector<uint8_t> %v_data = %v.data();", s, s),
		fmt.Sprintf("writer.append_bytes(%v_data);", s),
		fmt.Sprintf("writer.write_field_size(%d, %v_data.size());", field.Number, s))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("const std::vector<uint8_t> %v_data = %v.data();", varName, varName),
		fmt.Sprintf("writer.append_bytes(%v_data);", varName))
}
