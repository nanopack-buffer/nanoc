package tsgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type stringGenerator struct{}

func (g stringGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "string"
}

func (g stringGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%vByteLength + 4", varName)
}

func (g stringGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g stringGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g stringGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g stringGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)

	var cast string
	if field.Type.Kind == datatype.Enum {
		cast = fmt.Sprintf(" as T%v", field.Type.Identifier)
	}

	var l1 string
	if ctx.IsVariableInScope(c) {
		l1 = fmt.Sprintf("%v = reader.readString(ptr, %vByteLength)%v;", c, c, cast)
	} else {
		l1 = fmt.Sprintf("const %v = reader.readString(ptr, %vByteLength)%v;", c, c, cast)
	}

	return generator.Lines(
		fmt.Sprintf("const %vByteLength = reader.readFieldSize(%d);", c, field.Number),
		l1,
		fmt.Sprintf("ptr += %vByteLength;", c),
	)
}

func (g stringGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var cast string
	if dataType.Kind == datatype.Enum {
		cast = fmt.Sprintf(" as T%v", dataType.Identifier)
	}

	var l2 string
	if ctx.IsVariableInScope(varName) {
		l2 = fmt.Sprintf("%v = reader.readString(ptr, %vByteLength)%v;", varName, varName, cast)
	} else {
		l2 = fmt.Sprintf("const %v = reader.readString(ptr, %vByteLength)%v;", varName, varName, cast)
	}

	return generator.Lines(
		fmt.Sprintf("const %vByteLength = reader.readInt32(ptr);", varName),
		fmt.Sprintf("ptr + 4;"),
		l2,
		fmt.Sprintf("ptr += %vByteLength;"), varName)
}

func (g stringGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("const %vByteLength = writer.appendString(this.%v);", c, c),
		fmt.Sprintf("writer.writeFieldSize(%d, %vByteLength, offset);", field.Number, c),
		fmt.Sprintf("bytesWritten += %vByteLength;", c))
}

func (g stringGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("const %vByteLength = writer.appendStringAndSize(%v);", varName, varName)
}
