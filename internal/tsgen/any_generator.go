package tsgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type anyGenerator struct{}

func (g anyGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return "NanoPackReader"
}

func (g anyGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%v.byteLength", varName)
}

func (g anyGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g anyGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g anyGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g anyGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)

	var l1 string
	if ctx.IsVariableInScope(c) {
		l1 = fmt.Sprintf("%v = reader.newReaderAt(ptr, ptr + %vByteLength);", c, c)
	} else {
		l1 = fmt.Sprintf("const %v = reader.newReaderAt(ptr, ptr + %vByteLength);", c, c)
	}

	return generator.Lines(
		fmt.Sprintf("const %vByteLength = reader.readFieldSize(%d);", c, field.Number),
		l1,
		fmt.Sprintf("ptr += %vByteLength;", c))
}

func (g anyGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var l2 string
	if ctx.IsVariableInScope(varName) {
		l2 = fmt.Sprintf("%v = reader.newReaderAt(ptr, ptr + %vByteLength);", varName, varName)
	} else {
		l2 = fmt.Sprintf("const %v = reader.newReaderAt(ptr, ptr + %vByteLength);", varName, varName)
	}

	return generator.Lines(
		fmt.Sprintf("const %vByteLength = reader.readInt32(ptr);", varName),
		"ptr += 4",
		l2,
		fmt.Sprintf("ptr += %vByteLength;", varName))
}

func (g anyGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("writer.writeFieldSize(%d, this.%v.byteLength);", field.Number, c),
		fmt.Sprintf("writer.appendBytes(this.%v)", c))
}

func (g anyGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("writer.appendInt32(%v.byteLength);", varName),
		fmt.Sprintf("writer.appendBytes(%v);", varName))
}
