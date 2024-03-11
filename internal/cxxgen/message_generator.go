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
	if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		return fmt.Sprintf("std::unique_ptr<%v>", dataType.Identifier)
	}
	return dataType.Identifier
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v_byte_size", varName)
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	td := g.TypeDeclaration(field.Type)
	s := strcase.ToSnake(field.Name)
	if field.IsSelfReferencing() {
		return fmt.Sprintf("std::shared_ptr<%v> %v", field.Type.ElemType.Identifier, s)
	}
	return td + " " + s
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	if field.IsSelfReferencing() {
		return fmt.Sprintf("%v(std::move(%v))", s, s)
	}
	if ms := field.Type.Schema.(*npschema.Message); ms.IsInherited {
		return fmt.Sprintf("%v(std::move(%v))", s, s)
	}
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	if field.IsSelfReferencing() {
		return fmt.Sprintf("std::shared_ptr<%v> %v;", field.Type.ElemType.Identifier, s)
	}
	return g.TypeDeclaration(field.Type) + " " + s + ";"
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.IsSelfReferencing() {
		return generator.Lines(
			fmt.Sprintf("if (reader.read_field_size(%d) < 0) {", field.Number),
			fmt.Sprintf("    %v = nullptr;", s),
			"} else {",
			fmt.Sprintf("    int %v_bytes_read = 0;", s),
			fmt.Sprintf("    %v = std::make_shared<%v>(begin + ptr, %v_bytes_read);", s, field.Type.ElemType.Identifier, s),
			fmt.Sprintf("    ptr += %v_bytes_read;", s),
			"}")
	}

	td := g.TypeDeclaration(field.Type)
	ms := field.Type.Schema.(*npschema.Message)

	var ctor string
	if field.Type.Schema.(*npschema.Message).IsInherited {
		ctor = fmt.Sprintf("std::move(make_%v(begin + ptr, %v_bytes_read))", strcase.ToSnake(ms.Name), s)
	} else {
		ctor = fmt.Sprintf("%v(begin + ptr, %v_bytes_read)", td, s)
	}

	return generator.Lines(
		fmt.Sprintf("int %v_bytes_read = 0;", s),
		fmt.Sprintf("%v = %v;", s, ctor),
		fmt.Sprintf("ptr += %v_bytes_read;", s))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	td := g.TypeDeclaration(dataType)
	isPolymorphic := dataType.Schema.(*npschema.Message).IsInherited

	var expr string
	if isPolymorphic {
		expr = fmt.Sprintf("std::move(make_%v(begin + ptr, %v_bytes_read))", strcase.ToSnake(dataType.Identifier), varName)
	} else {
		expr = fmt.Sprintf("%v(begin + ptr, %v_bytes_read)", td, varName)
	}

	var l1 string
	if ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("%v = %v", varName, expr)
	} else if isPolymorphic {
		l1 = fmt.Sprintf("std::unique_ptr %v = %v;", varName, expr)
	} else {
		l1 = fmt.Sprintf("%v %v(begin + ptr, %v_bytes_read);", td, varName, varName)
	}

	return generator.Lines(
		fmt.Sprintf("int %v_bytes_read = 0;", varName),
		l1,
		fmt.Sprintf("ptr += %v_bytes_read;", varName))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.IsSelfReferencing() {
		return generator.Lines(
			fmt.Sprintf("if (%v != nullptr) {", s),
			fmt.Sprintf("    const size_t %v_byte_size = %v->write_to(buf, buf.size());", s, s),
			fmt.Sprintf("    NanoPack::write_field_size(%d, %v_byte_size, offset, buf);", field.Number, s),
			"} else {",
			fmt.Sprintf("    NanoPack::write_field_size(%d, -1, offset, buf);", field.Number),
			"}")
	}

	ms := field.Type.Schema.(*npschema.Message)

	var l0 string
	if ms.IsInherited {
		l0 = fmt.Sprintf("const size_t %v_byte_size = %v->write_to(buf, buf.size());", s, s)
	} else {
		l0 = fmt.Sprintf("const size_t %v_byte_size = %v.write_to(buf, buf.size());", s, s)
	}

	return generator.Lines(
		l0,
		fmt.Sprintf("NanoPack::write_field_size(%d, %v_byte_size, offset, buf);", field.Number, s))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ms := dataType.Schema.(*npschema.Message)
	if ms.IsInherited {
		return fmt.Sprintf("const size_t %v_byte_size = %v->write_to(buf, buf.size());", varName, varName)
	}
	return fmt.Sprintf("const size_t %v_byte_size = %v.write_to(buf, buf.size());", varName, varName)
}
