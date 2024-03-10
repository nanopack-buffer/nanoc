package tsgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strconv"
)

type numberGenerator struct{}

var appendMethodName = map[datatype.Kind]string{
	datatype.Int8:   "appendInt8",
	datatype.Int32:  "appendInt32",
	datatype.Int64:  "appendInt64",
	datatype.UInt8:  "appendUint8",
	datatype.UInt32: "appendUint32",
	datatype.UInt64: "appendUint64",
	datatype.Double: "appendDouble",
}

var readMethodName = map[datatype.Kind]string{
	datatype.Int8:   "readInt8",
	datatype.Int32:  "readInt32",
	datatype.Int64:  "readInt64",
	datatype.UInt8:  "readUint8",
	datatype.UInt32: "readUint32",
	datatype.UInt64: "readUint64",
	datatype.Double: "readDouble",
}

func (g numberGenerator) TypeDeclaration(dataType datatype.DataType) string {
	if dataType.Kind == datatype.Int64 {
		return "bigint"
	}
	return "number"
}

func (g numberGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return strconv.Itoa(dataType.ByteSize)
}

func (g numberGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g numberGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v", c, c)
}

func (g numberGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g numberGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return g.ReadValueFromBuffer(field.Type, strcase.ToLowerCamel(field.Name), ctx)
}

func (g numberGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var m string
	var cast string

	if dataType.Kind == datatype.Enum {
		cast = fmt.Sprintf(" as T%v", dataType.Identifier)
		m = readMethodName[dataType.ElemType.Kind]
	} else {
		m = readMethodName[dataType.Kind]
	}

	var l0 string
	if ctx.IsVariableInScope(varName) {
		l0 = fmt.Sprintf("%v = reader.%v(ptr)%v;", varName, m, cast)
	} else {
		l0 = fmt.Sprintf("const %v = reader.%v(ptr)%v;", varName, m, cast)
	}

	return generator.Lines(
		l0,
		fmt.Sprintf("ptr += %d;", dataType.ByteSize))
}

func (g numberGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.%v(this.%v);", appendMethodName[field.Type.Kind], strcase.ToLowerCamel(field.Name)),
		fmt.Sprintf("writer.writeFieldSize(%d, %d, offset);", field.Number, field.Type.ByteSize),
		fmt.Sprintf("bytesWritten += %d", field.Type.ByteSize))
}

func (g numberGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("writer.%v(%v);", appendMethodName[dataType.Kind], varName)
}
