package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type messageGenerator struct{}

func (g messageGenerator) TypeDeclaration(dataType datatype.DataType) string {
	if dataType.Schema == nil {
		return "std::unique_ptr<NanoPack::Message>"
	}
	if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		return fmt.Sprintf("std::unique_ptr<%v>", dataType.Identifier)
	}
	return dataType.Identifier
}

func (g messageGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	if dataType.Schema == nil {
		return "std::unique_ptr<NanoPack::Message> " + paramName
	}
	if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		return fmt.Sprintf("std::unique_ptr<%v> %v", dataType.Identifier, paramName)
	}
	return fmt.Sprintf("const %v &%v", dataType.Identifier, paramName)
}

func (g messageGenerator) RValue(dataType datatype.DataType, argName string) string {
	if dataType.Schema == nil {
		return fmt.Sprintf("std::move(%v)", argName)
	}
	if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		return fmt.Sprintf("std::move(%v)", argName)
	}
	return argName
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v_byte_size", varName)
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	td := g.TypeDeclaration(field.Type)
	s := strcase.ToSnake(field.Name)
	if field.Type.Schema == nil {
		return fmt.Sprintf("std::unique_ptr<NanoPack::Message> %v", s)
	}
	if field.IsSelfReferencing() {
		return fmt.Sprintf("std::unique_ptr<%v> %v", field.Type.ElemType.Identifier, s)
	}
	return td + " " + s
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	if field.Type.Schema == nil || field.IsSelfReferencing() {
		return fmt.Sprintf("%v(std::move(%v))", s, s)
	}
	if ms := field.Type.Schema.(*npschema.Message); ms.IsInherited {
		return fmt.Sprintf("%v(std::move(%v))", s, s)
	}
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	if field.Type.Schema == nil {
		return fmt.Sprintf("std::unique_ptr<NanoPack::Message> %v;", s)
	}
	if field.IsSelfReferencing() {
		return fmt.Sprintf("std::unique_ptr<%v> %v;", field.Type.ElemType.Identifier, s)
	}
	return g.TypeDeclaration(field.Type) + " " + s + ";"
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.Type.Schema == nil {
		return generator.Lines(
			"reader.buffer += ptr;",
			fmt.Sprintf("size_t %v_bytes_read;", s),
			fmt.Sprintf("%v = std::move(make_nanopack_message(reader, %v_bytes_read));", s, s),
			"reader.buffer = buf;",
			fmt.Sprintf("ptr += %v_bytes_read;", s))
	}

	if field.IsSelfReferencing() {
		return generator.Lines(
			fmt.Sprintf("if (reader.read_field_size(%d) < 0) {", field.Number),
			fmt.Sprintf("    %v = nullptr;", s),
			"} else {",
			fmt.Sprintf("    %v = std::make_unique<%v>();", s, field.Type.ElemType.Identifier),
			"reader.buffer += ptr;",
			fmt.Sprintf("    const size_t %v_bytes_read = %v->read_from(reader);", s, s),
			"reader.buffer = buf;",
			fmt.Sprintf("    ptr += %v_bytes_read;", s),
			"}")
	}

	ms := field.Type.Schema.(*npschema.Message)

	if field.Type.Schema.(*npschema.Message).IsInherited {
		return generator.Lines(
			fmt.Sprintf("%v = std::move(make_%v(reader));", s, strcase.ToSnake(ms.Name)),
			"reader.buffer += ptr;",
			fmt.Sprintf("const size_t %v_bytes_read = %v->read_from(reader);", s, s),
			"reader.buffer = buf;",
			fmt.Sprintf("ptr += %v_bytes_read;", s))
	}

	return generator.Lines(
		"reader.buffer += ptr;",
		fmt.Sprintf("const size_t %v_bytes_read = %v.read_from(reader);", s, s),
		"reader.buffer = buf;",
		fmt.Sprintf("ptr += %v_bytes_read;", s))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	td := g.TypeDeclaration(dataType)

	var declr string
	if !ctx.IsVariableInScope(varName) {
		declr = fmt.Sprintf("%v %v;", td, varName)
	}

	var read string
	if dataType.Schema == nil {
		read = generator.Lines(
			fmt.Sprintf("size_t %v_bytes_read;", varName),
			"reader.buffer += ptr;",
			fmt.Sprintf("%v = std::move(make_nanopack_message(reader, %v_bytes_read));", varName, varName),
			"reader.buffer = buf;")
	} else if dataType.Schema.(*npschema.Message).IsInherited {
		read = generator.Lines(
			fmt.Sprintf("size_t %v_bytes_read;", varName),
			"reader.buffer += ptr;",
			fmt.Sprintf("%v = std::move(make_%v(reader, %v_bytes_read));", varName, strcase.ToSnake(dataType.Identifier), varName),
			"reader.buffer = buf;")
	} else {
		read = generator.Lines(
			"reader.buffer += ptr;",
			fmt.Sprintf("const size_t %v_bytes_read = %v.read_from(reader);", varName, varName),
			"reader.buffer = buf;",
		)
	}

	return generator.Lines(
		declr,
		read,
		fmt.Sprintf("ptr += %v_bytes_read;", varName))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.IsSelfReferencing() {
		return generator.Lines(
			fmt.Sprintf("if (%v != nullptr) {", s),
			fmt.Sprintf("    const size_t %v_byte_size = %v->write_to(writer, writer.size());", s, s),
			fmt.Sprintf("    writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s),
			"} else {",
			fmt.Sprintf("    writer.write_field_size(%d, -1, offset);", field.Number),
			"}")
	}

	if field.Type.Schema == nil || field.Type.Kind == datatype.Optional {
		return generator.Lines(
			fmt.Sprintf("const size_t %v_byte_size = %v->write_to(writer, writer.size());", s, s),
			fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s))
	}

	ms := field.Type.Schema.(*npschema.Message)

	var sizeDeclr string
	if ms.IsInherited {
		sizeDeclr = fmt.Sprintf("const size_t %v_byte_size = %v->write_to(writer, writer.size());", s, s)
	} else {
		sizeDeclr = fmt.Sprintf("const size_t %v_byte_size = %v.write_to(writer, writer.size());", s, s)
	}

	return generator.Lines(
		sizeDeclr,
		fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if dataType.Schema == nil {
		return fmt.Sprintf("const size_t %v_byte_size = %v->write_to(writer, writer.size());", varName, varName)
	}
	if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		return fmt.Sprintf("const size_t %v_byte_size = %v->write_to(writer, writer.size());", varName, varName)
	}
	return fmt.Sprintf("const size_t %v_byte_size = %v.write_to(writer, writer.size());", varName, varName)
}
