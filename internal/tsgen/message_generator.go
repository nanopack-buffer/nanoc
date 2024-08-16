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
	if dataType.Schema == nil {
		return "NanoPackMessage"
	}
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
	c := strcase.ToLowerCamel(field.Name)
	tmpv := "maybe" + strcase.ToCamel(field.Name)

	var ctor string
	if field.Type.Schema == nil {
		ctor = "makeNanoPackMessage(reader, ptr)"
	} else if field.Schema.IsInherited {
		ctor = fmt.Sprintf("make%v(reader, ptr)", field.Type.Identifier)
	} else {
		ctor = fmt.Sprintf("%v.fromReader(reader, ptr)", g.TypeDeclaration(field.Type))
	}

	return generator.Lines(
		fmt.Sprintf("const %v = %v;", tmpv, ctor),
		fmt.Sprintf("if (!%v) { return null; }", tmpv),
		fmt.Sprintf("const %v = %v.result;", c, tmpv),
		fmt.Sprintf("ptr += %v.bytesRead;", tmpv))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	tmpv := "maybe" + strcase.ToCamel(varName)

	var ctor string
	if dataType.Schema == nil {
		ctor = "makeNanoPackMessage(reader, ptr)"
	} else if ms, ok := dataType.Schema.(*npschema.Message); ok && ms.IsInherited {
		ctor = fmt.Sprintf("make%v(reader, ptr)", dataType.Identifier)
	} else {
		ctor = fmt.Sprintf("%v.fromReader(reader, ptr)", g.TypeDeclaration(dataType))
	}

	var assignment string
	if ctx.IsVariableInScope(varName) {
		assignment = fmt.Sprintf("%v = %v.result;", varName, tmpv)
	} else {
		assignment = fmt.Sprintf("const %v = %v.result;", varName, tmpv)
	}

	return generator.Lines(
		fmt.Sprintf("const %v = %v;", tmpv, ctor),
		fmt.Sprintf("if (!%v) {", tmpv),
		"    return null;",
		"}",
		assignment,
		fmt.Sprintf("ptr += %v.bytesRead;", tmpv))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("const %vWriteOffset = writer.currentSize;", c),
		fmt.Sprintf("writer.reserveHeader(this.%v.headerSize);", c),
		fmt.Sprintf("const %vByteSize = this.%v.writeTo(writer, %vWriteOffset);", c, c, c),
		fmt.Sprintf("writer.writeFieldSize(%d, %vByteSize, offset);", field.Number, c),
		fmt.Sprintf("bytesWritten += %vByteSize;", c))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return generator.Lines(
		fmt.Sprintf("const %vWriteOffset = writer.currentSize;", varName),
		fmt.Sprintf("writer.reserveHeader(%v.headerSize);", varName),
		fmt.Sprintf("const %vByteSize = %v.writeTo(writer, %vWriteOffset);", varName, varName, varName))
}
