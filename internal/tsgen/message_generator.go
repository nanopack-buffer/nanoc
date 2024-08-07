package tsgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type messageGenerator struct{}

func (g messageGenerator) TypeDeclaration(dataType datatype.DataType) string {
	return dataType.Identifier
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return varName + "ByteSize"
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	return g.ReadValueFromBuffer(field.Type, strcase.ToLowerCamel(field.Name), ctx)
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	tmpv := "maybe" + strcase.ToCamel(varName)

	var l4 string
	if ctx.IsVariableInScope(varName) {
		l4 = fmt.Sprintf("%v = %v.result;", varName, tmpv)
	} else {
		l4 = fmt.Sprintf("const %v = %v.result;", varName, tmpv)
	}

	return generator.Lines(
		fmt.Sprintf("const %v = %v.fromReader(reader, ptr);", tmpv, g.TypeDeclaration(dataType)),
		fmt.Sprintf("if (!%v) {", tmpv),
		"    return null;",
		"}",
		l4,
		fmt.Sprintf("ptr += %v.bytesRead;", tmpv))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		"const offset = writer.currentSize;",
		fmt.Sprintf("writer.reserveHeader(%v.headerSize);", c),
		fmt.Sprintf("const %vByteSize = this.%v.writeTo(writer, offset);", c, c),
		fmt.Sprintf("writer.writeFieldSize(%d, %vByteSize, offset);", field.Number, c),
		fmt.Sprintf("bytesWritten += %vByteSize;", c))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		"const offset = writer.currentSize;",
		fmt.Sprintf("writer.reserveHeader(%v.headerSize);", varName),
		fmt.Sprintf("const %vByteSize = %v.writeTo(writer, offset);", varName, varName))
}
