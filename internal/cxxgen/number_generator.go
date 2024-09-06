package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

var cxxIntTypes = map[datatype.Kind]string{
	datatype.Int8:   "int8_t",
	datatype.Int32:  "int32_t",
	datatype.Int64:  "int64_t",
	datatype.Double: "double",
	datatype.UInt8:  "uint8_t",
	datatype.UInt32: "uint32_t",
	datatype.UInt64: "uint64_t",
}

// numberGenerator is a DataTypeMessageCodeGenerator for NanoPack number types.
type numberGenerator struct{}

func (g numberGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return cxxIntTypes[dataType.Kind]
}

func (g numberGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return cxxIntTypes[dataType.Kind] + " " + paramName
}

func (g numberGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g numberGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%d", dataType.ByteSize)
}

func (g numberGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return cxxIntTypes[field.Type.Kind] + " " + strcase.ToSnake(field.Name)
}

func (g numberGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(%v)", s, s)
}

func (g numberGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g numberGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	switch field.Type.Kind {
	case datatype.Enum:
		return g.ReadValueFromBuffer(*field.Type.ElemType, s+"_raw_value", ctx)

	case datatype.Optional:
		return generator.Lines(
			fmt.Sprintf("%v %v_value;", cxxIntTypes[field.Type.ElemType.Kind], s),
			fmt.Sprintf("reader.read_%v(ptr, %v_value);", field.Type.ElemType.Identifier, s),
			fmt.Sprintf("%v = %v_value;", s, s),
			fmt.Sprintf("ptr += %d;", field.Type.ElemType.ByteSize))

	default:
		return generator.Lines(
			fmt.Sprintf("reader.read_%v(ptr, %v);", field.Type.Identifier, s),
			fmt.Sprintf("ptr += %d;", field.Type.ByteSize))
	}
}

func (g numberGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	if dataType.Kind == datatype.Optional {
		return generator.Lines(
			fmt.Sprintf("%v %v_value;", cxxIntTypes[dataType.ElemType.Kind], varName),
			fmt.Sprintf("reader.read_%v(ptr, %v_value);", dataType.Identifier, varName),
			fmt.Sprintf("%v = %v_value;", varName, varName),
			fmt.Sprintf("ptr += %d;", dataType.ByteSize))
	}

	var declr string
	if !ctx.IsVariableInScope(varName) {
		declr = fmt.Sprintf("%v %v;", cxxIntTypes[dataType.Kind], varName)
	}

	return generator.Lines(
		declr,
		fmt.Sprintf("reader.read_%v(ptr, %v);", dataType.Identifier, varName),
		fmt.Sprintf("ptr += %d;", dataType.ByteSize))
}

func (g numberGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	switch field.Type.Kind {
	case datatype.Enum:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %d, offset);", field.Number, field.Type.ElemType.ByteSize),
			fmt.Sprintf("writer.append_%v(%v.value());", field.Type.ElemType.Identifier, strcase.ToSnake(field.Name)))
	case datatype.Optional:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %d, offset);", field.Number, field.Type.ElemType.ByteSize),
			fmt.Sprintf("writer.append_%v(*%v);", field.Type.ElemType.Identifier, strcase.ToSnake(field.Name)))

	default:
		return generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %d, offset);", field.Number, field.Type.ByteSize),
			fmt.Sprintf("writer.append_%v(%v);", field.Type.Identifier, strcase.ToSnake(field.Name)))
	}
}

func (g numberGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	switch dataType.Kind {
	case datatype.Enum:
		return fmt.Sprintf("writer.append_%v(%v.value());", dataType.ElemType.Identifier, varName)

	case datatype.Optional:
		return fmt.Sprintf("writer.append_%v(*%v);", dataType.ElemType.Identifier, varName)

	default:
		return fmt.Sprintf("writer.append_%v(%v);", dataType.Identifier, varName)
	}
}
